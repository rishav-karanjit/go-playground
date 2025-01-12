package main

import (
	"fmt"
)

func someTask(id int, data chan int) {
	for taskID := range data {
		fmt.Printf("Worker: %d executed Task %d\n", id, taskID)
	}
}

func main() {
	// Creating a channel
	channel := make(chan int)

	// Creating 2 workers to execute the task
	for i := 0; i < 2; i++ {
		go someTask(i, channel)
	}
	// Filling channel with 5 numbers to be executed
	for i := 0; i < 5; i++ {
		channel <- i
	}
}
