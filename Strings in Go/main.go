package main

import (
	"fmt"
)

func main() {
	s := "Hellö 世界"
	incorrectUse(s)
	correctUse(s)

	surrogatePair := "😊"
	incorrectUse(surrogatePair)
	correctUse(surrogatePair)
}

// Incorrectly using the string directly in Go
func incorrectUse(goString string) {
	fmt.Println("Incorrect Use:")
	// Iterate over bytes in Go
	for i := 0; i < len(goString); i++ {
		fmt.Printf("%c ", goString[i])
	}
	fmt.Println("\n")
}

// correct use
func correctUse(goString string) {
	fmt.Println("Correct Use:")
	goStringToBytes := []byte(goString)
	goStringCorrected := string(goStringToBytes)
	runes := []rune(goStringCorrected)
	for i := 0; i < len(runes); i++ {
		fmt.Printf("%c", runes[i])
	}
	fmt.Println("\n\n")
}

// runes allows to handle text in a way that respects the Unicode standard
// runes is int32
