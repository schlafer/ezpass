package ezpass

import (
	"math"
	"unicode"
)

// CalculateEntropy returns the key-space entropy in bits based on length and pool size.
// If poolSize <= 0, it dynamically estimates the pool size based on character classes present.
func CalculateEntropy(password string, poolSize int) float64 {
	if len(password) == 0 {
		return 0
	}

	if poolSize <= 0 {
		poolSize = estimatePoolSize(password)
	}

	if poolSize <= 1 {
		return 0
	}

	runes := []rune(password)
	return float64(len(runes)) * math.Log2(float64(poolSize))
}

// EstimatePoolSize estimates the effective character set size of a given string.
func estimatePoolSize(password string) int {
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	specialCount := 0
	seenSpecial := make(map[rune]bool)

	for _, r := range password {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			if !seenSpecial[r] {
				seenSpecial[r] = true
				specialCount++
			}
			hasSpecial = true
		}
	}

	pool := 0
	if hasLower {
		pool += 26
	}
	if hasUpper {
		pool += 26
	}
	if hasDigit {
		pool += 10
	}
	if hasSpecial {
		if specialCount > 6 {
			pool += 32
		} else {
			pool += specialCount
		}
	}

	if pool == 0 {
		pool = 1
	}
	return pool
}

// StrengthRating returns a human-readable strength description based on bits of entropy.
func StrengthRating(entropy float64) string {
	switch {
	case entropy < 36:
		return "Very Weak"
	case entropy < 60:
		return "Weak"
	case entropy < 80:
		return "Fair"
	case entropy < 110:
		return "Strong"
	default:
		return "Very Strong"
	}
}

// CalculatePassphraseEntropy returns entropy for a passphrase based on word count and wordlist size.
func CalculatePassphraseEntropy(wordCount, wordlistSize int) float64 {
	if wordCount <= 0 || wordlistSize <= 1 {
		return 0
	}
	return float64(wordCount) * math.Log2(float64(wordlistSize))
}

// ShannonEntropy calculates the information entropy (uncertainty) of a string.
func ShannonEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	counts := make(map[rune]float64)
	runes := []rune(s)
	total := float64(len(runes))

	for _, r := range runes {
		counts[r]++
	}

	var entropy float64
	for _, count := range counts {
		p := count / total
		entropy -= p * math.Log2(p)
	}
	return entropy
}
