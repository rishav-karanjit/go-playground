package main

import (
	"fmt"
	"reflect"
)

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

	switch v := union.(type) {

	case *MyUnionMemberStringValue:
		fmt.Println("Union holds a string", v.Value)
	default:
		panic(fmt.Sprintf("Union holds an unknown type %T ", v))
	}
	// fmt.Println("Union holds an integer %v", union)
	fmt.Println(reflect.TypeOf(union))
	// union = &MyUnionMemberStringValue{
	// 	Value: "Hello",
	// }
	// fmt.Println("Union holds a string %v", union)
}
