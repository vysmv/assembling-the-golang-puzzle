package main

import "fmt"

func dangerous() {
    fmt.Println("before panic!") // Отобразится первым
    panic("boom") // завершаю эту функцию и запускаю defer если таковые есть. 
	fmt.Println("after panic") 
}

func foo() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recovered in foo:", r) // отобразится вторым recovered in foo: boom
        }
    }()

    dangerous() // frame #3
    fmt.Println("after dangerous") // НЕ ВЫПОЛНЯЕТСЯ!!! Так как тдет раскрутка стека и при возварет из dangerous сразу происхдит завершение foo и вызов ее defer.
}

func main() { // frame #1
    foo() // frame #2
    fmt.Println("after foo") // отобразится третим 
}

