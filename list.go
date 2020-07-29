package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sort"
)

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
			if StringInSlice(v, *values) == false {
				*values = append(*values, v)
			}
		}
		return nil
	}
}

func List(key string) {
	if key == "" {
		uniqueKeys := make([]string, 0)
		w := walk{prefix: prefix, abstractArray: &uniqueKeys}
		w.WalkFrontMatter(w.WalkListKeys)
		sort.Strings(uniqueKeys)
		PrintArrayOnNewlines(uniqueKeys)
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
