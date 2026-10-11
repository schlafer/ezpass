package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/schlafer/ezpass/pkg/ezpass"
)

const Version = "1.1.0"

// Run executes the ezpass CLI application.
func Run() {
	var (
		length           int
		count            int
		noUpper          bool
		noLower          bool
		noDigits         bool
		noSymbols        bool
		onlyAlphanumeric bool
		symbols          string
		excludeAmbiguous bool
		excludeChars     string

		passphraseMode bool
		wordsCount     int
		separator      string
		capitalize     bool
		includeNum     bool

		copyToClipboard bool
		showEntropy     bool
		noNewline       bool
		showVersion     bool
	)

	fs := flag.NewFlagSet("ezpass", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	// Password flags (short & long)
	fs.IntVar(&length, "length", ezpass.DefaultLength, "Password length")
	fs.IntVar(&length, "l", ezpass.DefaultLength, "Password length (shorthand)")

	fs.IntVar(&count, "count", 1, "Number of passwords to generate")
	fs.IntVar(&count, "n", 1, "Number of passwords to generate (shorthand)")

	fs.BoolVar(&noUpper, "no-upper", false, "Exclude uppercase letters (A-Z)")
	fs.BoolVar(&noUpper, "U", false, "Exclude uppercase letters (shorthand)")

	fs.BoolVar(&noLower, "no-lower", false, "Exclude lowercase letters (a-z)")
	fs.BoolVar(&noLower, "L", false, "Exclude lowercase letters (shorthand)")

	fs.BoolVar(&noDigits, "no-digits", false, "Exclude numbers (0-9)")
	fs.BoolVar(&noDigits, "D", false, "Exclude numbers (shorthand)")

	fs.BoolVar(&noSymbols, "no-symbols", false, "Exclude special symbols")
	fs.BoolVar(&noSymbols, "S", false, "Exclude special symbols (shorthand)")

	fs.BoolVar(&onlyAlphanumeric, "only-alphanumeric", false, "Use only letters and digits (shorthand -A)")
	fs.BoolVar(&onlyAlphanumeric, "A", false, "Use only letters and digits")

	fs.StringVar(&symbols, "symbols", "", "Custom symbols (default: !@#$*_)")

	fs.BoolVar(&excludeAmbiguous, "exclude-ambiguous", false, "Exclude lookalike characters (0, O, 1, l, I, |)")
	fs.BoolVar(&excludeAmbiguous, "a", false, "Exclude lookalike characters (shorthand)")

	fs.StringVar(&excludeChars, "exclude", "", "Exclude specific characters")
	fs.StringVar(&excludeChars, "e", "", "Exclude specific characters (shorthand)")

	// Passphrase flags
	fs.BoolVar(&passphraseMode, "passphrase", false, "Generate a memorable Diceware passphrase")
	fs.BoolVar(&passphraseMode, "p", false, "Generate a memorable passphrase (shorthand)")

	fs.IntVar(&wordsCount, "words", ezpass.DefaultPassphraseWordCount, "Number of words in passphrase")
	fs.IntVar(&wordsCount, "w", ezpass.DefaultPassphraseWordCount, "Number of words in passphrase (shorthand)")

	fs.StringVar(&separator, "separator", "-", "Word separator for passphrase")
	fs.StringVar(&separator, "sep", "-", "Word separator for passphrase (shorthand)")

	fs.BoolVar(&capitalize, "capitalize", false, "Capitalize each word in passphrase")
	fs.BoolVar(&includeNum, "include-number", false, "Append a random number to passphrase")

	// Output & utility flags
	fs.BoolVar(&copyToClipboard, "copy", false, "Copy generated password to clipboard")
	fs.BoolVar(&copyToClipboard, "c", false, "Copy to clipboard (shorthand)")

	fs.BoolVar(&showEntropy, "entropy", false, "Display entropy and strength rating")
	fs.BoolVar(&showEntropy, "E", false, "Display entropy and strength (shorthand)")

	fs.BoolVar(&noNewline, "no-newline", false, "Do not print trailing newline")
	fs.BoolVar(&noNewline, "0", false, "Do not print trailing newline (shorthand)")

	fs.BoolVar(&showVersion, "version", false, "Show version information")
	fs.BoolVar(&showVersion, "v", false, "Show version information (shorthand)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `ezpass - Secure CLI password & passphrase generator (v%s)

USAGE:
  ezpass [flags]

EXAMPLES:
  ezpass                             # Generate a default 12-char secure password
  ezpass -l 20                       # Generate a 20-char password
  ezpass -n 5                        # Generate 5 passwords
  ezpass -l 16 -a                    # Exclude ambiguous characters (0, O, 1, l, I, |)
  ezpass -l 16 -S                    # Password with no special symbols
  ezpass -c                          # Generate and copy to clipboard
  ezpass -p -w 4                     # Generate 4-word passphrase (e.g. correct-horse-battery-staple)
  ezpass -p -w 4 --capitalize --include-number

FLAGS:
  -l, --length <int>                 Password length (default: 12)
  -n, --count <int>                  Number of passwords to generate (default: 1)
  -S, --no-symbols                   Exclude special characters
  -D, --no-digits                    Exclude numeric digits
  -U, --no-upper                     Exclude uppercase letters
  -L, --no-lower                     Exclude lowercase letters
  -A, --only-alphanumeric            Use only letters and numbers
      --symbols <str>                Custom symbol set (default: "!@#$*_")
  -a, --exclude-ambiguous            Exclude lookalike characters (0, O, 1, l, I, |)
  -e, --exclude <str>                Exclude specific characters

PASSPHRASE OPTIONS:
  -p, --passphrase                   Generate a memorable passphrase
  -w, --words <int>                  Word count for passphrase (default: 4)
      --separator, --sep <str>       Word separator (default: "-")
      --capitalize                   Capitalize each word
      --include-number               Append a random number to the passphrase

OUTPUT & UTILITY:
  -c, --copy                         Copy output to clipboard
  -E, --entropy                      Display bits of entropy and strength analysis
  -0, --no-newline                   Omit trailing newline (convenient for scripts)
  -v, --version                      Show ezpass version
  -h, --help                         Show help message
`, Version)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(2)
	}

	if showVersion {
		fmt.Printf("ezpass version %s\n", Version)
		os.Exit(0)
	}

	if count <= 0 {
		fmt.Fprintln(os.Stderr, "Error: count must be greater than 0")
		os.Exit(1)
	}

	cfg := ezpass.Config{
		Length:           length,
		NoUpper:          noUpper,
		NoLower:          noLower,
		NoDigits:         noDigits,
		NoSymbols:        noSymbols,
		OnlyAlphanumeric: onlyAlphanumeric,
		Symbols:          symbols,
		ExcludeAmbiguous: excludeAmbiguous,
		ExcludeChars:     excludeChars,
	}

	passCfg := ezpass.PassphraseConfig{
		WordCount:  wordsCount,
		Separator:  separator,
		Capitalize: capitalize,
		IncludeNum: includeNum,
	}

	var results []string

	for i := 0; i < count; i++ {
		var (
			output   string
			poolSize int
			err      error
			entropy  float64
		)

		if passphraseMode {
			output, poolSize, err = ezpass.GeneratePassphrase(passCfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			entropy = ezpass.CalculatePassphraseEntropy(wordsCount, poolSize)
		} else {
			output, poolSize, err = ezpass.Generate(cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			entropy = ezpass.CalculateEntropy(output, poolSize)
		}

		results = append(results, output)

		if showEntropy {
			rating := ezpass.StrengthRating(entropy)
			if count == 1 {
				fmt.Printf("Password:  %s\nEntropy:   %.1f bits (%s)\nPool Size: %d\n", output, entropy, rating, poolSize)
			} else {
				fmt.Printf("%-24s [Entropy: %.1f bits, Strength: %s]\n", output, entropy, rating)
			}
		}
	}

	if !showEntropy {
		if noNewline && count == 1 {
			fmt.Print(results[0])
		} else {
			for _, res := range results {
				fmt.Println(res)
			}
		}
	}

	if copyToClipboard {
		clipboardText := strings.Join(results, "\n")
		if err := ezpass.CopyToClipboard(clipboardText); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, "✓ Copied to clipboard")
		}
	}
}
