package main

import "fmt"

func main() {
	defer fmt.Println("defer in main")
	b()
}

func a() {
	defer fmt.Println("defer in a")
	panic("something went wrong")
}

func b() {
	defer fmt.Println("defer in b")
	a()
}