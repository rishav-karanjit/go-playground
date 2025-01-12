package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	stuff := []uint8{240, 168, 137, 159, 240, 168, 137, 159, 240, 168, 137, 159}
	byteSlice := []byte(stuff)
	str := string(byteSlice)

	// r, _ := utf8.DecodeRuneInString(str)
	fmt.Print(utf8.ValidString(str))

}
