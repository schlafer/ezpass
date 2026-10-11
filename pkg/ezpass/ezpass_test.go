package ezpass

import (
	"strings"
	"testing"
	"unicode"
)

func TestGenerate_Defaults(t *testing.T) {
	cfg := DefaultConfig()
	pass, poolSize, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pass) != DefaultLength {
		t.Fatalf("expected length %d, got %d (%s)", DefaultLength, len(pass), pass)
	}
	if poolSize <= 0 {
		t.Fatalf("expected poolSize > 0, got %d", poolSize)
	}

	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, r := range pass {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case strings.ContainsRune(DefaultSymbols, r):
			hasSymbol = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSymbol {
		t.Errorf("default password should contain characters from all 4 groups, got: %s", pass)
	}
}

func TestGenerate_VariousLengths(t *testing.T) {
	lengths := []int{1, 2, 3, 4, 8, 16, 32, 64, 128}
	for _, l := range lengths {
		cfg := Config{Length: l}
		pass, _, err := Generate(cfg)
		if err != nil {
			t.Fatalf("failed at length %d: %v", l, err)
		}
		if len([]rune(pass)) != l {
			t.Errorf("expected length %d, got %d", l, len([]rune(pass)))
		}
	}
}

func TestGenerate_InvalidLength(t *testing.T) {
	_, _, err := Generate(Config{Length: 0})
	if err == nil {
		t.Errorf("expected error for length 0, got nil")
	}

	_, _, err = Generate(Config{Length: -5})
	if err == nil {
		t.Errorf("expected error for negative length, got nil")
	}
}

func TestGenerate_NoSymbols(t *testing.T) {
	cfg := Config{
		Length:    50,
		NoSymbols: true,
	}
	pass, _, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range pass {
		if strings.ContainsRune(DefaultSymbols, r) {
			t.Errorf("found symbol %c when NoSymbols was requested: %s", r, pass)
		}
	}
}

func TestGenerate_OnlyAlphanumeric(t *testing.T) {
	cfg := Config{
		Length:           50,
		OnlyAlphanumeric: true,
	}
	pass, _, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range pass {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			t.Errorf("found non-alphanumeric character %c in password: %s", r, pass)
		}
	}
}

func TestGenerate_NoUpper(t *testing.T) {
	cfg := Config{
		Length:  50,
		NoUpper: true,
	}
	pass, _, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range pass {
		if unicode.IsUpper(r) {
			t.Errorf("found uppercase letter %c when NoUpper was set: %s", r, pass)
		}
	}
}

func TestGenerate_NoLower(t *testing.T) {
	cfg := Config{
		Length:  50,
		NoLower: true,
	}
	pass, _, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range pass {
		if unicode.IsLower(r) {
			t.Errorf("found lowercase letter %c when NoLower was set: %s", r, pass)
		}
	}
}

func TestGenerate_NoDigits(t *testing.T) {
	cfg := Config{
		Length:   50,
		NoDigits: true,
	}
	pass, _, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range pass {
		if unicode.IsDigit(r) {
			t.Errorf("found digit %c when NoDigits was set: %s", r, pass)
		}
	}
}

func TestGenerate_CustomSymbols(t *testing.T) {
	custom := "@#~"
	cfg := Config{
		Length:  50,
		Symbols: custom,
	}
	pass, _, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range pass {
		if strings.ContainsRune("!$*_", r) {
			t.Errorf("found default symbol %c not in custom set %s", r, custom)
		}
	}
}

func TestGenerate_ExcludeAmbiguous(t *testing.T) {
	cfg := Config{
		Length:           100,
		ExcludeAmbiguous: true,
	}
	pass, _, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range pass {
		if strings.ContainsRune(AmbiguousChars, r) {
			t.Errorf("found ambiguous character %c: %s", r, pass)
		}
	}
}

func TestGenerate_ExcludeChars(t *testing.T) {
	forbidden := "abcABC123"
	cfg := Config{
		Length:       100,
		ExcludeChars: forbidden,
	}
	pass, _, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range pass {
		if strings.ContainsRune(forbidden, r) {
			t.Errorf("found excluded character %c: %s", r, pass)
		}
	}
}

func TestGenerate_EmptyPool(t *testing.T) {
	cfg := Config{
		Length:    12,
		NoUpper:   true,
		NoLower:   true,
		NoDigits:  true,
		NoSymbols: true,
	}
	_, _, err := Generate(cfg)
	if err == nil {
		t.Fatalf("expected error for empty pool, got nil")
	}
}

func TestGeneratePassphrase(t *testing.T) {
	cfg := PassphraseConfig{
		WordCount:  4,
		Separator:  "-",
		Capitalize: false,
	}
	pass, poolSize, err := GeneratePassphrase(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if poolSize != len(DefaultWordList) {
		t.Errorf("expected pool size %d, got %d", len(DefaultWordList), poolSize)
	}
	parts := strings.Split(pass, "-")
	if len(parts) != 4 {
		t.Fatalf("expected 4 parts, got %d (%s)", len(parts), pass)
	}

	// Test with capitalization and number
	cfg2 := PassphraseConfig{
		WordCount:  3,
		Separator:  ".",
		Capitalize: true,
		IncludeNum: true,
	}
	pass2, _, err := GeneratePassphrase(cfg2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parts2 := strings.Split(pass2, ".")
	if len(parts2) != 4 { // 3 words + 1 number
		t.Fatalf("expected 4 parts with number, got %d (%s)", len(parts2), pass2)
	}
	for i := 0; i < 3; i++ {
		firstRune := rune(parts2[i][0])
		if !unicode.IsUpper(firstRune) {
			t.Errorf("word %s should be capitalized", parts2[i])
		}
	}
}

func TestEntropyAndStrength(t *testing.T) {
	weak := CalculateEntropy("abc", 26)
	if weak >= 36 {
		t.Errorf("expected weak entropy, got %f", weak)
	}
	if rating := StrengthRating(weak); rating != "Very Weak" {
		t.Errorf("expected Very Weak, got %s", rating)
	}

	strong := CalculateEntropy("A!b@C#d$E%f^1&2*", 70)
	if strong < 80 {
		t.Errorf("expected strong entropy, got %f", strong)
	}
	if rating := StrengthRating(strong); rating != "Strong" && rating != "Very Strong" {
		t.Errorf("expected Strong or Very Strong, got %s", rating)
	}
}

func TestZeroRunesAndBytes(t *testing.T) {
	r := []rune("secretPassword")
	ZeroRunes(r)
	for _, val := range r {
		if val != 0 {
			t.Errorf("rune was not zeroed: %v", val)
		}
	}

	b := []byte("secretPassword")
	ZeroBytes(b)
	for _, val := range b {
		if val != 0 {
			t.Errorf("byte was not zeroed: %v", val)
		}
	}
}

func BenchmarkGeneratePassword(b *testing.B) {
	cfg := DefaultConfig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = Generate(cfg)
	}
}

func BenchmarkGeneratePassphrase(b *testing.B) {
	cfg := PassphraseConfig{
		WordCount: 4,
		Separator: "-",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = GeneratePassphrase(cfg)
	}
}
