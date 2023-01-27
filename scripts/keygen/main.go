package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(GenerateKey(32))
}

// GenerateKey Generates alpha-numeric key of 'size' parameter
func GenerateKey(size int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	key := make([]rune, size)
	for index := range key {
		key[index] = letters[rand.Intn(len(letters))]
	}

	return string(key)
}
