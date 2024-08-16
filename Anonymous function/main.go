package main

import "fmt"

func main() {
	value := *func() *string {
		a := "hello"
		return &a
	}()
	fmt.Println(value)

	value2 := func() *string {
		a := "hello"
		return &a
	}()
	fmt.Println(value2)
}
