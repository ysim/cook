package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var invalidFileCount = 0

func ValidateSingleFile(fullFilepath string) error {
	recipeFile, err := ParseFile(fullFilepath)
	if err != nil {
		fmt.Printf("%s (%s)\n", fullFilepath, err.Error())
		invalidFileCount++
		return nil
	}

	_, err = ParseFrontMatter(recipeFile.FrontMatter)
	if err != nil {
		fmt.Printf("%s (%s)\n", fullFilepath, err.Error())
		invalidFileCount++
		return nil
	}
	return nil
}

func ValidateFile() filepath.WalkFunc {
	return func(fullFilepath string, info os.FileInfo, err error) error {
		shouldSkip := ShouldSkipFile(info, err)
		if shouldSkip {
			return nil
		}

		return ValidateSingleFile(fullFilepath)
	}
}

// Scan the list of files in `prefix` and output all the ones that don't
// conform to the formatting.
func ValidateFiles() {
	err := filepath.Walk(prefix, ValidateFile())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("# invalid files: %d\n", invalidFileCount)
}
