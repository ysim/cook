package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	orOperator = ","
)

func SearchFile(args map[string][]string) filepath.WalkFunc {
	return func(fullFilepath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		// Don't descend into directories for now
		if info.IsDir() {
			return nil
		}
		// Ignore hidden files
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		recipeFile, err := ParseFile(fullFilepath)
		if err != nil {
			return nil
		}
		frontMatter, err := ParseFrontMatter(recipeFile.FrontMatter)
		if err != nil {
			log.Printf("Unknown type detected in front matter of file '%s'\n", fullFilepath)
		}

		var showFile bool // zero value is false
		for argKey, argValueArray := range args {
			// ok is set to true if the key exists, false if not
			fileValueArray, ok := frontMatter[argKey]

			if !ok {
				return nil
			} else {
				for _, argValue := range argValueArray {
					for _, fileValue := range fileValueArray {
						if strings.Contains(fileValue, argValue) {
							showFile = true
							break
						}
					}
				}
			}
		}
		if showFile {
			fmt.Println(GetBasenameWithoutExt(fullFilepath))
		}
		return nil
	}
}

func ParseSearchQuery(args []string) (map[string][]string, error) {
	q := make(map[string][]string)

	// Consolidate all arguments
	argString := strings.Join(args[:], " ")

	// Now split into keys
	fields := strings.Fields(argString)

	for _, f := range fields {
		// strings.Split will always return an array of at least one item
		// (if there are no matches, that item will be an empty string)
		splitField := strings.Split(f, ":")
		if len(splitField) != 2 {
			return nil, errors.New("Exactly one ':' is required per whitespace-delimited argument")
		}

		key, value := splitField[0], splitField[1]
		valueArray := strings.Split(value, orOperator)
		q[key] = valueArray
	}
	return q, nil
}

func Search(args []string) {
	parsedQuery, parseErr := ParseSearchQuery(args)
	if parseErr != nil {
		fmt.Printf("Invalid search query: %s\n", parseErr.Error())
		os.Exit(1)
	}

	searchErr := filepath.Walk(prefix, SearchFile(parsedQuery))
	if searchErr != nil {
		log.Fatal(searchErr)
	}
}
