package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	a := string(uint8(255))
	fmt.Println(a)
	fmt.Println(utf8.ValidString(a))

	// Create the byte slice
	invalidBytes := []byte{uint8(255)}
	stringinvalidBytes := string(invalidBytes)
	// Check if the byte slice is a valid UTF-8 sequence
	if utf8.ValidString(stringinvalidBytes) {
		fmt.Println("The byte sequence is valid UTF-8")
	} else {
		fmt.Println("The byte sequence is not valid UTF-8")
	}
}
