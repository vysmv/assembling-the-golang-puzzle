package main

import (
	"fmt"
	"os"
)

func main() {
	defer fmt.Println("defer executed")

	fmt.Println("before return")

	if len(os.Args) > 1 && os.Args[1] == "return" {
		return
	}

	fmt.Println("after return")
}