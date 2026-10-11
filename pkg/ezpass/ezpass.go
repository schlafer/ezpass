package ezpass

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"unicode"
)

const (
	DefaultLength              = 12
	DefaultSymbols             = "!@#$*_"
	AmbiguousChars             = "0O1lI|"
	DefaultPassphraseWordCount = 4
	DefaultPassphraseSeparator = "-"
)

// Config defines the options for generating a random password.
type Config struct {
	Length           int
	NoUpper          bool
	NoLower          bool
	NoDigits         bool
	NoSymbols        bool
	OnlyAlphanumeric bool
	Symbols          string
	ExcludeAmbiguous bool
	ExcludeChars     string
}

// PassphraseConfig defines options for generating a memorable passphrase.
type PassphraseConfig struct {
	WordCount  int
	Separator  string
	Capitalize bool
	IncludeNum bool
	WordList   []string
}

// DefaultConfig returns standard recommended defaults.
func DefaultConfig() Config {
	return Config{
		Length:  DefaultLength,
		Symbols: DefaultSymbols,
	}
}

// SecureRandom generates a cryptographically secure random integer in [min, max].
func SecureRandom(min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("min (%d) cannot be greater than max (%d)", min, max)
	}
	diff := int64(max - min + 1)
	nBig, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return 0, fmt.Errorf("crypto/rand failure: %w", err)
	}
	return int(nBig.Int64()) + min, nil
}

// filterGroup removes excluded characters and optional ambiguous characters from a charset.
func filterGroup(charset string, excludeAmbiguous bool, excludeChars string) string {
	var b strings.Builder
	excludeSet := make(map[rune]bool)

	for _, r := range excludeChars {
		excludeSet[r] = true
	}
	if excludeAmbiguous {
		for _, r := range AmbiguousChars {
			excludeSet[r] = true
		}
	}

	for _, r := range charset {
		if !excludeSet[r] {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Generate creates a cryptographically secure password based on Config.
// Returns the generated password, the total pool size used, and any error encountered.
func Generate(cfg Config) (string, int, error) {
	if cfg.Length <= 0 {
		return "", 0, errors.New("password length must be greater than 0")
	}

	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	digitChars := "0123456789"
	symbolChars := cfg.Symbols
	if symbolChars == "" && !cfg.NoSymbols && !cfg.OnlyAlphanumeric {
		symbolChars = DefaultSymbols
	}

	var activeGroups []string

	if !cfg.NoUpper {
		filtered := filterGroup(upperChars, cfg.ExcludeAmbiguous, cfg.ExcludeChars)
		if len(filtered) > 0 {
			activeGroups = append(activeGroups, filtered)
		}
	}

	if !cfg.NoLower {
		filtered := filterGroup(lowerChars, cfg.ExcludeAmbiguous, cfg.ExcludeChars)
		if len(filtered) > 0 {
			activeGroups = append(activeGroups, filtered)
		}
	}

	if !cfg.NoDigits {
		filtered := filterGroup(digitChars, cfg.ExcludeAmbiguous, cfg.ExcludeChars)
		if len(filtered) > 0 {
			activeGroups = append(activeGroups, filtered)
		}
	}

	if !cfg.NoSymbols && !cfg.OnlyAlphanumeric {
		filtered := filterGroup(symbolChars, cfg.ExcludeAmbiguous, cfg.ExcludeChars)
		if len(filtered) > 0 {
			activeGroups = append(activeGroups, filtered)
		}
	}

	if len(activeGroups) == 0 {
		return "", 0, errors.New("character pool is empty; verify your inclusion and exclusion options")
	}

	totalPoolSize := 0
	for _, g := range activeGroups {
		totalPoolSize += len([]rune(g))
	}

	numGroups := len(activeGroups)
	passRunes := make([]rune, 0, cfg.Length)

	if cfg.Length >= numGroups {
		base := cfg.Length / numGroups
		remainder := cfg.Length % numGroups

		for g := 0; g < numGroups; g++ {
			count := base
			if g < remainder {
				count++
			}
			groupRunes := []rune(activeGroups[g])
			for i := 0; i < count; i++ {
				idx, err := SecureRandom(0, len(groupRunes)-1)
				if err != nil {
					return "", 0, err
				}
				passRunes = append(passRunes, groupRunes[idx])
			}
		}
	} else {
		// Length is smaller than number of active groups.
		// Randomly select distinct groups to pick from.
		perm, err := randPerm(numGroups, cfg.Length)
		if err != nil {
			return "", 0, err
		}
		for _, g := range perm {
			groupRunes := []rune(activeGroups[g])
			idx, err := SecureRandom(0, len(groupRunes)-1)
			if err != nil {
				return "", 0, err
			}
			passRunes = append(passRunes, groupRunes[idx])
		}
	}

	if err := Shuffle(passRunes); err != nil {
		return "", 0, err
	}

	return string(passRunes), totalPoolSize, nil
}

// randPerm generates a random permutation of k distinct integers chosen from [0, n).
func randPerm(n, k int) ([]int, error) {
	if k > n {
		return nil, fmt.Errorf("cannot select %d distinct elements from %d", k, n)
	}
	pool := make([]int, n)
	for i := range pool {
		pool[i] = i
	}
	for i := 0; i < k; i++ {
		j, err := SecureRandom(i, n-1)
		if err != nil {
			return nil, err
		}
		pool[i], pool[j] = pool[j], pool[i]
	}
	return pool[:k], nil
}

// Shuffle performs a cryptographically secure in-place Fisher–Yates shuffle.
func Shuffle(runes []rune) error {
	for i := len(runes) - 1; i > 0; i-- {
		j, err := SecureRandom(0, i)
		if err != nil {
			return err
		}
		runes[i], runes[j] = runes[j], runes[i]
	}
	return nil
}

// GeneratePassphrase produces a memorable Diceware-style passphrase.
// Returns the passphrase, the wordlist size, and any error encountered.
func GeneratePassphrase(cfg PassphraseConfig) (string, int, error) {
	wordList := cfg.WordList
	if len(wordList) == 0 {
		wordList = DefaultWordList
	}
	if len(wordList) == 0 {
		return "", 0, errors.New("wordlist is empty")
	}

	wordsCount := cfg.WordCount
	if wordsCount <= 0 {
		wordsCount = DefaultPassphraseWordCount
	}

	sep := cfg.Separator
	if sep == "" && cfg.Separator != "" {
		sep = DefaultPassphraseSeparator
	}

	selectedWords := make([]string, wordsCount)
	for i := 0; i < wordsCount; i++ {
		idx, err := SecureRandom(0, len(wordList)-1)
		if err != nil {
			return "", 0, err
		}
		w := wordList[idx]
		if cfg.Capitalize && len(w) > 0 {
			runes := []rune(w)
			runes[0] = unicode.ToUpper(runes[0])
			w = string(runes)
		}
		selectedWords[i] = w
	}

	result := strings.Join(selectedWords, sep)

	if cfg.IncludeNum {
		num, err := SecureRandom(0, 99)
		if err != nil {
			return "", 0, err
		}
		result = fmt.Sprintf("%s%s%d", result, sep, num)
	}

	return result, len(wordList), nil
}

// ZeroRunes clears a slice of runes from memory.
func ZeroRunes(r []rune) {
	for i := range r {
		r[i] = 0
	}
}

// ZeroBytes clears a slice of bytes from memory.
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
