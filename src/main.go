package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
)

func main() {
	length := flag.Int("length", 12, "Password length")
	flag.Parse()

	if *length <= 0 {
		fmt.Fprintln(os.Stderr, "Length must be greater than 0")
		os.Exit(1)
	}

	password := generatePassword(*length)
	shuffled := []rune(password)
	shuffle(shuffled)
	fmt.Println(string(shuffled))
}

// --- Secure random number generator ---
func secureRandom(min, max int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		panic(err)
	}
	return int(nBig.Int64()) + min
}

// --- Generate password using 4 character groups ---
func generatePassword(length int) string {
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower := "abcdefghijklmnopqrstuvwxyz"
	nums := "0123456789"
	special := "!@#$*_"

	groups := []string{upper, lower, nums, special}
	groupSizes := []int{len(upper), len(lower), len(nums), len(special)}

	base := length / 4
	remainder := length % 4
	pass := make([]rune, 0, length)

	for g := 0; g < 4; g++ {
		thisLen := base
		if g < remainder {
			thisLen++
		}
		for i := 0; i < thisLen; i++ {
			idx := secureRandom(0, groupSizes[g]-1)
			pass = append(pass, rune(groups[g][idx]))
		}
	}

	return string(pass)
}

// --- Fisher–Yates shuffle ---
func shuffle(runes []rune) {
	for i := len(runes) - 1; i > 0; i-- {
		j := secureRandom(0, i)
		runes[i], runes[j] = runes[j], runes[i]
	}
}
