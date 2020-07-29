package main

import (
	"fmt"
	"sort"
)

func List(key string) {
	if key == "" {
		uniqueKeys := make([]string, 0)
		w := walk{prefix: prefix, abstractArray: &uniqueKeys}
		listErr := w.WalkFrontMatter(w.WalkListKeys)

		if listErr != nil {
			errMsg := "An error occurred while attempting to list keys."
			fmt.Printf(errMsg, key)
		}

		sort.Strings(uniqueKeys)
		PrintArrayOnNewlines(uniqueKeys)
	} else {
		uniqueValues := make([]string, 0)
		w := walk{prefix: prefix, abstractArray: &uniqueValues, key: key}
		listErr := w.WalkFrontMatter(w.ListValuesForKey)

		if listErr != nil {
			errMsg := "An error occurred while attempting to list values for key '%s'."
			fmt.Printf(errMsg, key)
		}

		if len(uniqueValues) == 0 {
			fmt.Printf("No values were found for key '%s'.\n", key)
			return
		}

		sort.Strings(uniqueValues)
		PrintArrayOnNewlines(uniqueValues)
	}
}
