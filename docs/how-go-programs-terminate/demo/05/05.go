package main

import (
	"fmt"
	"os"
)

func main() {
	defer fmt.Println("defer executed") // не будет вызвана

	fmt.Println("before exit")
	os.Exit(127)

	fmt.Println("after exit") // не будет вызвана
}