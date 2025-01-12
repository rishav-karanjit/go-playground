package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	// UTF-8 encoded string
	utf8Str := "Hello, 世界"
	fmt.Printf("UTF-8 string: %s\n", utf8Str)

	// Get the byte slice representation of the UTF-8 string
	utf8Bytes := []byte(utf8Str)
	fmt.Printf("UTF-8 bytes: %X\n", utf8Bytes)

	// Get the rune slice representation of the UTF-8 string
	runes := []rune(utf8Str)
	fmt.Printf("Runes: %c\n", runes)

	// Iterate over runes
	for i, r := range runes {
		fmt.Printf("%d: %c\n", i, r)
	}

	// UTF-16 encoded string
	utf16Str := "\U0001F600" // U+1F600 (Grinning Face emoji)
	fmt.Printf("UTF-16 string: %s\n", utf16Str)

	// Get the rune slice representation of the UTF-16 string
	runes16 := []rune(utf16Str)
	fmt.Printf("Runes (UTF-16): %c\n", runes16)

	// Check if a string is valid UTF-8
	invalidUTF8 := "\xc3\x28"
	if !utf8.ValidString(invalidUTF8) {
		fmt.Println("Invalid UTF-8 string")
	}
}
