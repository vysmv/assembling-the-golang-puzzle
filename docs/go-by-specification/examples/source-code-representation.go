package main

import "fmt"

func main() {
	s1 := "\u00e4"   // "ä" как одна кодовая точка (U+00E4)
	s2 := "a\u0308"  // "ä" как комбинация "a" (Перед обратным слешем) (U+0061) + "¨" (U+0308)

	s3 := "ä"
	s4 := "ä"

	s5 := "ä"
	s6 := "ä"
	
	fmt.Println(s1 == s2) // false
	fmt.Println(len(s1))  // 2 (байта в UTF-8)
	fmt.Println(len(s2))  // 3 (байта в UTF-8)
	fmt.Println(s1)  // ä
	fmt.Println(s2)  // ä

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

	/*If such situations arise, it is worth comparing the sizes of the strings in bytes.*/
	if len(s3) != len(s4) {
		fmt.Println("the number of bytes is different")
	}
}

