package main

import "fmt"

func main() {
	s1 := "\u00e4"   
	s2 := "a\u0308"


	s3 := "ä"
	s4 := "ä"

	s5 := "ä"
	s6 := "ä"
	
	fmt.Println(s1 == s2) 
	fmt.Println(len(s1))  
	fmt.Println(len(s2))  

	fmt.Println(s1) 
	fmt.Println(s2) 

	if s3 == s4 {
		fmt.Println("s3 and s4 are equal")
	} else {
		fmt.Println("s3 and s4 are not equal")
	}

	if s5 == s6 {
		fmt.Println("s5 and s6 are equal")
	} else {
		fmt.Println("s5 and s6 are not equal")
	}

	if len(s1) != len(s2) {
		fmt.Println("the number of bytes is different")
	}
}