# ezpass

CLI password generator in Go.

- Guarantees uppercase, lowercase, digits, and special symbols
- Avoids modulo bias via rejection sampling
- Shuffles output with Fisher–Yates algorithm

# Usage

- `go build main.go`
- `./main --length 20`
