package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func SearchFile(args map[string]string) filepath.WalkFunc {
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

		for argKey, argValue := range args {
			fileValue, ok := frontMatter[argKey]
			if !ok {
				// Key doesn't exist in the front matter.
				return nil
			} else {
				for _, v := range fileValue {
					if strings.Contains(v, argValue) {
						fmt.Println(GetBasenameWithoutExt(fullFilepath))
					}
				}
			}
		}
		return nil
	}
}

func Search(args []string) {
	// For now, just allow simple searching, for one value on one field, e.g.
	// "ingredients=chicken"
	argString := strings.Join(args[:], " ")
	keyValue := strings.Split(argString, ":")

	// strings.Split will always return an array of at least one item
	// (if there are no matches, that item will be an empty string)
	switch len(keyValue) {
	case 2:
		key := keyValue[0]
		value := keyValue[1]
		a := map[string]string{key: value}
		err := filepath.Walk(prefix, SearchFile(a))
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Invalid search query.")
	}
}
