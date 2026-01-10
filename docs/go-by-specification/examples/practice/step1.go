package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type State struct {
	Current int
	Total   int
}

func main() {

	path := getWorkDir()
	
	check(path)
	
	app, err := initApp(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//var input string

	for {

		content, err := getContent(path, app)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		break
		// Рендеринг. Очистить экран.

		// Напечатать рамку.

		// Напечатать содержание.

		// Напечатать под рамкой [1...n]=open a slide by its number / [n]=next / [p]=prev [q]=quit >

		// Считать ввод.

		// Обновить сосотояние. Записав новое значение текущего элемента.
		// По формуле
		// Если ввод - `n`:
		// state.Current < state.Total то state.Current++ иначе state.Current = 1
		// Если ввод - `p`:
		// state.Current > 1 то state.Current-- иначе state.Current = Total
		// Если ввод - число от 1...n:
		// Проверяем наличие такого слайда и если есть то
		// state.Current = указанное число
		// Если нет то печатаем slide {numder} not found
		// Если ввод - `q`:
		// break
		// Если ввод любой другой символ:
		// continue

		// Переход на следующую итерацию.
	}
}

func getWorkDir() string {

	var path string

	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		fmt.Println("The directory with the slides " +
			"is not specified")
		os.Exit(1)
	}

	return path

}

/*
Checks:
- that the ./slides directory exists
- that there is at least one file named "1"
- that slide filenames are numbered sequentially with a step of 1 (no gaps)
*/

func check(path string) {
	/*
		Note

		1. The os package is used to work with filesystem objects
			(files, directories) and to perform standard operations
			such as:
			- CRUD,
			- permissions,
			- ownership,
			- etc.

		2. The path/filepath package is used to work with filesystem
			paths as strings.
	*/

	info, err := os.Stat(path)

	if err != nil {
		fmt.Println("[check]\n" +
			"The directory with the slides is not specified")
		os.Exit(2)
	}

	if !info.IsDir() {
		fmt.Println("[check]\n" +
			path + ": is not dir")
		os.Exit(2)
	}

	entries, err := os.ReadDir(path)

	if err != nil {
		fmt.Println("[check]\n" + err.Error())
		os.Exit(2)
	}

	var nums []int

	for _, e := range entries {
		if e.IsDir() {
			fmt.Println("[check]\n" + "The directory [" + path + "]" +
				" must contain only text files from 1 to n")
			os.Exit(2)
		}

		n, convErr := strconv.Atoi(e.Name())
		if convErr != nil {
			fmt.Println("[check]\n" + "invalid slide filename %q " +
				"(must be a number)")
			os.Exit(2)
		}

		if n <= 0 {
			fmt.Printf("[check]\n"+"invalid slide number %d\n", n)
			os.Exit(2)
		}

		nums = append(nums, n)

		if len(nums) == 0 {

		}

		sort.Ints(nums)

		if nums[0] != 1 {
			fmt.Println("[check]\n" + "slide 1 not found")
			os.Exit(2)
		}

		for i := 1; i < len(nums); i++ {
			if nums[i] != nums[i-1]+1 {
				fmt.Printf("slides must be numbered "+
					" 1..N without gaps; expected %d, got %d",
					nums[i-1]+1,
					nums[i])

				os.Exit(2)
			}
		}

	}

}

/*
Initialize the application’s initial state:

	State{
		Current: 1,
		Total:   number of files in the directory,
	}
*/

func initApp(path string) (State, error) {

	/*
		A function is:
		“Take input data → do the work → return a result or error information.”
	*/
	allEntry, err := os.ReadDir(path)

	if err != nil {
		return State{}, errors.New("Failed to read directory " +
			"during initialization")
	}

	total := len(allEntry)
	return State{
			Current: 1,
			Total:   total,
		},
		nil
}

func getContent(path string, app State) (string, error) {

	currentName := strconv.Itoa(app.Current)
	filePath := filepath.Join(path, currentName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil

}
