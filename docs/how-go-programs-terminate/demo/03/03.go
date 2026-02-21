package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		fmt.Println("goroutine started")
		panic("boom in goroutine")
	}()

	time.Sleep(time.Second)
	fmt.Println("main finished")
}