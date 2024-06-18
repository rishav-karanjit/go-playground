package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	a := string(uint8(255))
	fmt.Println(a)
	fmt.Println(utf8.ValidString(a))
}
