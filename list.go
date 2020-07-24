package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sort"
)

func stringInSlice(s string, slice []string) bool {
	for _, i := range slice {
		if s == i {
			return true
		}
	}
	return false
}

func listValuesForKey(key string, values *[]string) filepath.WalkFunc {
	return func(fullFilepath string, info os.FileInfo, err error) error {
		shouldSkip := ShouldSkipFile(info, err)
		if shouldSkip {
			return nil
		}

		recipeFile, err := ParseFile(fullFilepath)
		if err != nil {
			log.WithFields(log.Fields{
				"file": fullFilepath,
			}).Warn(err.Error())
			return nil
		}
		frontMatter, err := ParseFrontMatter(recipeFile.FrontMatter)
		if err != nil {
			log.WithFields(log.Fields{
				"file": fullFilepath,
			}).Warn("Unknown type detected in front matter")
		}

		valueOfKey := frontMatter[key]
		for _, v := range valueOfKey {
			if stringInSlice(v, *values) == false {
				*values = append(*values, v)
			}
		}
		return nil
	}
}

func List(key string) {
	if key == "" {
		fmt.Println("Functionality to list all keys is not implemented yet. (#21)")
	} else {
		values := make([]string, 0)
		listErr := filepath.Walk(prefix, listValuesForKey(key, &values))
		if len(values) == 0 {
			fmt.Printf("No values were found for key '%s'.\n", key)
			return
		}
		sort.Strings(values)
		PrintArrayOnNewlines(values)
		if listErr != nil {
			errMsg := "An error occurred while attempting to list values for key '%s'."
			fmt.Printf(errMsg, key)
		}
	}
}
