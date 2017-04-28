package main

import (
	"fmt"
	"log"
	"os"
	"path"
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
		cleanedFilename := strings.Replace(path.Base(fullFilepath), path.Ext(fullFilepath), "", -1)
		splitBytesArray := ParseFile(cleanedFilename)
		frontMatter := ParseFrontMatter(splitBytesArray[1])

		for argKey, argValue := range args {
			fileValue, ok := frontMatter[argKey]
			if !ok {
				// Key doesn't exist in the front matter.
				return nil
			} else {
				// Type assertions
				fmt.Println(fileValue)
				// TODO: errors when searching on string fields
				fileValueCoerced := fileValue.([]interface{})
				fileValueArray := make([]string, len(fileValueCoerced))
				for _, v := range fileValueCoerced {
					fileValueArray = append(fileValueArray, v.(string))
				}

				// Now check for the argValue
				for _, v := range fileValueArray {
					if strings.Contains(v, argValue) {
						fmt.Println(cleanedFilename)
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
	keyValue := strings.Split(argString, "=")

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
