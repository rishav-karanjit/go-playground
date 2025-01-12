package main

import "fmt"

func pointerIp(x *int) {
	fmt.Println(*x)
}

func main() {
	input := 10
	pointerIp(&input)
}
