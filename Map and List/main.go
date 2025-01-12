package main

import "fmt"

func main() {
	// Create a map
	myMap := map[string]int{
		"apple":  5,
		"banana": 3,
		"orange": 7,
	}

	// Iterate over the map
	fmt.Println("Iterating over the map:")
	for key, value := range myMap {
		fmt.Printf("Key: %s, Value: %d\n", key, value)
	}

	// Create a list (slice)
	myList := []string{"hello", "world", "go", "programming"}

	// Iterate over the list
	fmt.Println("\nIterating over the list:")
	for index, value := range myList {
		fmt.Printf("Index: %d, Value: %s\n", index, value)
	}

	emptyList := []string{}
	fmt.Println(emptyList == nil)
	fmt.Println(emptyList)
}
