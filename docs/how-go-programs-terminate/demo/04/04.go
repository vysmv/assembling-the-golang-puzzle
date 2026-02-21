package main

import (
	"fmt"
	"time"
)

func main() {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered:", r)
			}
		}()

		panic("boom")
	}()

	time.Sleep(time.Second)
	fmt.Println("main finished")
}