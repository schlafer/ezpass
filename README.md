# ezpass

[![CI](https://github.com/schlafer/ezpass/actions/workflows/ci.yml/badge.svg)](https://github.com/schlafer/ezpass/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A fast, cryptographically secure CLI password and memorable passphrase generator written in Go.

---

## Features

- **Cryptographically Secure**: Powered by Go's `crypto/rand` for high-entropy randomness.
- **Fair Character Distribution**: Guarantees characters across enabled groups and shuffles using Fisher–Yates.
- **Memorable Passphrases**: Diceware-style passphrases using the official EFF wordlist (`--passphrase`).
- **Granular Character Controls**: Disable symbols, numbers, uppercase, or lowercase; or restrict to alphanumeric only.
- **Exclude Ambiguous Characters**: Filter out confusing lookalikes like `0`, `O`, `1`, `l`, `I`, and `|`.
- **Clipboard Integration**: Seamlessly copy passwords via `-c` / `--copy` with zero external dependencies (supports OSC 52, `wl-copy`, `xclip`, `pbcopy`).
- **Strength & Entropy Analysis**: Real-time keyspace and Shannon entropy calculation with human-readable ratings (`--entropy`).
- **Batch Generation**: Generate multiple passwords at once (`-n 5`).
- **Script-Friendly**: Suppress trailing newlines (`-0`) for clean piping into commands, variables, or files.
- **Zero External Dependencies**: Pure Go standard library.
- **Reusable Go Package**: Can be imported into any Go application via `github.com/schlafer/ezpass/pkg/ezpass`.

---

## Installation

### Via `go install` (Recommended)

```bash
go install github.com/schlafer/ezpass@latest
```

### From Source

Clone the repository and build:

```bash
git clone https://github.com/schlafer/ezpass.git
cd ezpass
go build -o ezpass .
```

Optionally move it to your system PATH:

```bash
sudo mv ezpass /usr/local/bin/
```

---

## Quick Start & Usage

### 1. Standard Passwords

Generate a default 12-character password:

```bash
ezpass
# Output: Aq#8Tz3LpR_4
```

Specify a custom length:

```bash
ezpass -l 24
# or: ezpass --length 24
```

Generate a batch of 5 passwords:

```bash
ezpass -n 5
```

### 2. Character Filtering

Exclude special symbols:

```bash
ezpass -l 16 -S
# or: ezpass -l 16 --no-symbols
```

Alphanumeric only (letters and numbers):

```bash
ezpass -l 16 -A
# or: ezpass -l 16 --only-alphanumeric
```

Exclude ambiguous/lookalike characters (`0`, `O`, `1`, `l`, `I`, `|`):

```bash
ezpass -l 16 -a
# or: ezpass -l 16 --exclude-ambiguous
```

Exclude specific characters (e.g. quotes or slashes):

```bash
ezpass -l 16 -e "\"'/\\"
```

Use a custom symbol set:

```bash
ezpass -l 16 --symbols "!@#$%"
```

### 3. Memorable Passphrases (Diceware Style)

Generate a memorable 4-word passphrase:

```bash
ezpass -p
# Output: sleet-bony-cried-curse
```

Customize word count, capitalization, and numbers:

```bash
ezpass -p -w 5 --sep "_" --capitalize --include-number
# Output: Sleet_Bony_Cried_Curse_Frost_42
```

### 4. Clipboard & Entropy

Generate and copy directly to system clipboard:

```bash
ezpass -c
# ✓ Copied to clipboard
```

Display password entropy and strength analysis:

```bash
ezpass -l 20 -E
# Output:
# Password:  Ksx5@R2O5iP_zE@a05_*
# Entropy:   121.7 bits (Very Strong)
# Pool Size: 68
```

---

## Shell & Scripting Integration

Assign to a shell variable without trailing newline:

```bash
PASSWORD=$(ezpass -l 16 -0)
```

Pipe directly into `sha256sum` or files:

```bash
ezpass -l 32 -0 > secret.key
```

Batch generation into a CSV/text file:

```bash
ezpass -n 20 > passwords.txt
```

---

## CLI Flag Reference

| Flag | Shorthand | Description | Default |
| --- | --- | --- | --- |
| `--length <int>` | `-l` | Length of generated password | `12` |
| `--count <int>` | `-n` | Number of passwords to generate | `1` |
| `--no-symbols` | `-S` | Exclude special characters (`!@#$*_`) | `false` |
| `--no-digits` | `-D` | Exclude numeric digits (`0-9`) | `false` |
| `--no-upper` | `-U` | Exclude uppercase letters (`A-Z`) | `false` |
| `--no-lower` | `-L` | Exclude lowercase letters (`a-z`) | `false` |
| `--only-alphanumeric` | `-A` | Use only letters and numbers | `false` |
| `--symbols <str>` | | Custom symbol set | `"!@#$*_"` |
| `--exclude-ambiguous` | `-a` | Exclude lookalike characters (`0O1lI\|`) | `false` |
| `--exclude <str>` | `-e` | Exclude specific characters | `""` |
| `--passphrase` | `-p` | Generate memorable passphrase | `false` |
| `--words <int>` | `-w` | Word count for passphrase | `4` |
| `--separator <str>` | `--sep` | Word separator for passphrase | `"-"` |
| `--capitalize` | | Capitalize words in passphrase | `false` |
| `--include-number` | | Append random number to passphrase | `false` |
| `--copy` | `-c` | Copy result to system clipboard | `false` |
| `--entropy` | `-E` | Display entropy and strength | `false` |
| `--no-newline` | `-0` | Suppress trailing newline | `false` |
| `--version` | `-v` | Display version information | |
| `--help` | `-h` | Display help screen | |

---

## Using as a Go Library

You can import `ezpass` directly into your Go projects:

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
 fmt.Println("Generated:", password)
 fmt.Printf("Entropy: %.1f bits (%s)\n",
  ezpass.CalculateEntropy(password, poolSize),
  ezpass.StrengthRating(ezpass.CalculateEntropy(password, poolSize)),
 )

 // Memorable passphrase
 passCfg := ezpass.PassphraseConfig{
  WordCount:  4,
  Separator:  "-",
  Capitalize: true,
 }
 passphrase, _, _ := ezpass.GeneratePassphrase(passCfg)
 fmt.Println("Passphrase:", passphrase)
}
```

---

## Testing & Benchmarks

Run all unit tests:

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
