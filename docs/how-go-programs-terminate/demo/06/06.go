package main

import "fmt"

func dangerous() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in dangerous:", r)
		}
	}()

	fmt.Println("before panic")
	panic("boom")
	fmt.Println("after panic") 
}

func main() {
	fmt.Println("calling dangerous")
	dangerous()
	fmt.Println("after dangerous")
}