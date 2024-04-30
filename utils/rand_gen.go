package utils

import (
	"math/rand"
)

type StrGenFunc func() string

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	numberOfStrings = 8
)

func RandomStringGenerator() StrGenFunc {
	return func() string {
		b := make([]rune, numberOfStrings)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}
}
