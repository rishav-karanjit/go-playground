package main

import "fmt"

type MyUnion interface {
	isMyUnion()
}

type MyUnionMemberIntegerValue struct {
	Value int32
}

// MyUnionMemberIntegerValue implements MyUnion by giving implementing isMyUnion
func (*MyUnionMemberIntegerValue) isMyUnion() {}

type MyUnionMemberStringValue struct {
	Value string
}

func (*MyUnionMemberStringValue) isMyUnion() {}

func main() {
	var union MyUnion

	union = &MyUnionMemberIntegerValue{
		Value: 123,
	}
	fmt.Println("Union holds an integer %v", union)

	union = &MyUnionMemberStringValue{
		Value: "Hello",
	}
	fmt.Println("Union holds a string %v", union)
}
