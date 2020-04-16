package main

import (
	"fmt"
	"os"
)

func validateNewRecipe(filename string) {
	filepath := GetFullFilepath(filename)
	fmt.Println(filepath)
	_, err := os.Stat(filepath)

	// Not the same as os.IsExist! This is because os.Stat doesn't throw an
	// error if the file exists, so os.IsExist would receive a nil value for err.
	// So, a os.IsExist(err) block would never execute.
	if !os.IsNotExist(err) {
		fmt.Printf("There already exists a file at the path: %s\n", filepath)
		os.Exit(1)
	}
}

func CreateNewRecipe(filename string, name string) {
	validateNewRecipe(filename)
}
