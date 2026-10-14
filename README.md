# ezpass

[![CI](https://github.com/schlafer/ezpass/actions/workflows/ci.yml/badge.svg)](https://github.com/schlafer/ezpass/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/schlafer/ezpass.svg)](https://pkg.go.dev/github.com/schlafer/ezpass)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A command-line password and passphrase generator written in Go with zero external dependencies.

---

## Features

- **CSPRNG Generation**: Uses Go's `crypto/rand` for cryptographically secure randomness.
- **Fair Character Distribution**: Enforces representation from every active character class and shuffles with Fisher–Yates.
- **Diceware Passphrases**: Generates memorable passphrases using the EFF 7,776-word long wordlist.
- **Character Filtering**: Toggle symbols, digits, uppercase, or lowercase; restrict to alphanumeric; exclude specific characters or ambiguous lookalikes (`0`, `O`, `1`, `l`, `I`, `|`).
- **Clipboard Integration**: Copies passwords via system utilities (`wl-copy`, `xclip`, `xsel`, `pbcopy`, `clip.exe`) and ANSI OSC 52 escape sequences (works across SSH and tmux).
- **Entropy Analysis**: Calculates keyspace and Shannon entropy with strength ratings.
- **Script-Friendly**: Supports batch generation (`-n`) and trailing newline suppression (`-0`) for direct shell variable assignment and piping.
- **Zero Dependencies**: Pure Go standard library. Can also be imported as a library (`pkg/ezpass`).

---

## Installation

### Via `go install` (Recommended)

Requires Go 1.22+:

```bash
go install github.com/schlafer/ezpass@latest
```

Ensure `$GOPATH/bin` (or `$(go env GOPATH)/bin`) is included in your `PATH`.

### From Source

```bash
git clone https://github.com/schlafer/ezpass.git
cd ezpass
go build -o ezpass .
sudo mv ezpass /usr/local/bin/
```

---

## Quick Start & Usage

### 1. Passwords

Generate a default 12-character password:

```bash
ezpass
# Output: Aq#8Tz3LpR_4
```

Specify custom length:

```bash
ezpass -l 24
```

Generate multiple passwords:

```bash
ezpass -n 5
```

### 2. Character Filtering

Exclude special symbols:

```bash
ezpass -l 16 -S
```

Alphanumeric only (letters and digits):

```bash
ezpass -l 16 -A
```

Exclude ambiguous/lookalike characters (`0`, `O`, `1`, `l`, `I`, `|`):

```bash
ezpass -l 16 -a
```

Exclude specific characters (e.g., quotes, slashes):

```bash
ezpass -l 16 -e "\"'/\\"
```

Use a custom symbol set:

```bash
ezpass -l 16 --symbols "!@#$%"
```

Exclude digits, uppercase, or lowercase:

```bash
ezpass -D   # no numbers (0-9)
ezpass -U   # no uppercase (A-Z)
ezpass -L   # no lowercase (a-z)
```

### 3. Memorable Passphrases (Diceware Style)

Generate a 4-word passphrase:

```bash
ezpass -p
# Output: sleet-bony-cried-curse
```

Customize word count, separator, capitalization, and appended number:

```bash
ezpass -p -w 5 --sep "_" --capitalize --include-number
# Output: Sleet_Bony_Cried_Curse_Frost_42
```

### 4. Clipboard & Entropy

Generate and copy directly to the system clipboard:

```bash
ezpass -c
# ✓ Copied to clipboard
```

Display entropy and strength rating:

```bash
ezpass -l 20 -E
# Output:
# Password:  Ksx5@R2O5iP_zE@a05_*
# Entropy:   121.7 bits (Very Strong)
# Pool Size: 68
```

### 5. Shell & Scripting Integration

Store directly in a shell variable without trailing newline:

```bash
PASSWORD=$(ezpass -l 16 -0)
```

Pipe directly into files or commands:

```bash
ezpass -l 32 -0 > secret.key
ezpass -n 20 > passwords.txt
```

---

## CLI Flag Reference

| Flag | Shorthand | Description | Default |
| --- | --- | --- | --- |
| `--length <int>` | `-l` | Password length | `12` |
| `--count <int>` | `-n` | Number of passwords to generate | `1` |
| `--no-symbols` | `-S` | Exclude special characters (`!@#$*_`) | `false` |
| `--no-digits` | `-D` | Exclude numeric digits (`0-9`) | `false` |
| `--no-upper` | `-U` | Exclude uppercase letters (`A-Z`) | `false` |
| `--no-lower` | `-L` | Exclude lowercase letters (`a-z`) | `false` |
| `--only-alphanumeric` | `-A` | Use only letters and digits | `false` |
| `--symbols <str>` | | Custom symbol set | `"!@#$*_"` |
| `--exclude-ambiguous` | `-a` | Exclude lookalike characters (`0O1lI\|`) | `false` |
| `--exclude <str>` | `-e` | Exclude specific characters | `""` |
| `--passphrase` | `-p` | Generate memorable passphrase | `false` |
| `--words <int>` | `-w` | Word count for passphrase | `4` |
| `--separator <str>` | `--sep` | Word separator for passphrase | `"-"` |
| `--capitalize` | | Capitalize each word in passphrase | `false` |
| `--include-number` | | Append random number to passphrase | `false` |
| `--copy` | `-c` | Copy result to system clipboard | `false` |
| `--entropy` | `-E` | Display entropy and strength analysis | `false` |
| `--no-newline` | `-0` | Suppress trailing newline | `false` |
| `--version` | `-v` | Show version information | |
| `--help` | `-h` | Show help screen | |

---

## Using as a Go Library

Import `pkg/ezpass` directly in Go code:

```go
package main

import (
	"fmt"
	"log"

	"github.com/schlafer/ezpass/pkg/ezpass"
)

func main() {
	// Standard password
	cfg := ezpass.Config{
		Length:           20,
		ExcludeAmbiguous: true,
	}
	password, poolSize, err := ezpass.Generate(cfg)
	if err != nil {
		log.Fatal(err)
	}

	entropy := ezpass.CalculateEntropy(password, poolSize)
	fmt.Printf("Password: %s (%.1f bits, %s)\n", password, entropy, ezpass.StrengthRating(entropy))

	// Memorable passphrase
	passCfg := ezpass.PassphraseConfig{
		WordCount:  4,
		Separator:  "-",
		Capitalize: true,
	}
	passphrase, _, err := ezpass.GeneratePassphrase(passCfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Passphrase:", passphrase)
}
```

---

## Testing & Benchmarks

Run unit tests:

```bash
go test -v ./...
```

Run benchmarks:

```bash
go test -bench=. ./pkg/ezpass
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
