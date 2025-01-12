package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	emojiString := "\U0001F600" // U+1F600 (Grinning Face emoji)

	fmt.Printf("UTF-16 encoding: % X\n", []rune(emojiString))
	fmt.Printf("UTF-8 encoding: % X\n", []byte(emojiString))

	fmt.Printf("Length in UTF-16: %d\n", len([]rune(emojiString)))
	fmt.Printf("Length in UTF-8: %d\n", utf8.RuneCountInString(emojiString))
}
