# ezpass

CLI password generator in Go.

## Installation

Clone the repository:

```bash
git clone https://github.com/schlafer/ezpass
cd ezpass/src
```

Build the binary:

```bash
go build -o ezpass main.go
```

You can now move it somewhere in your PATH (optional):

```bash
sudo mv ezpass /usr/local/bin/
```

## Usage

Generate a 12-character password:

```bash
ezpass
```


Generate a 20-character password:

```bash
ezpass --length 20
```

Example output:

```bash
Aq#8Tz3LpR_4fKm1Vb
```

## Example Integration

Use it in a shell script:

```bash
PASSWORD=$(ezpass --length 16)
echo "Generated password: $PASSWORD"
```

Or redirect it to a file:

```bash
ezpass --length 24 > mypassword.txt
```

## How It Works

- Password is built from four character groups:

  - Uppercase letters (A–Z)

  - Lowercase letters (a–z)

  - Digits (0–9)

  - Special symbols (!@#$*_)

- Each group contributes equally to the total length.

- Final password is shuffled using the Fisher–Yates algorithm.

- Random numbers are generated via Go’s crypto/rand package for maximum entropy.
