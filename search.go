package main

import (
	"fmt"
	"strings"
)

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
		a := map[string]interface{}{
			key: []interface{}{value},
		}
		return a
	default:
		fmt.Println("Invalid search query.")
	}
}
