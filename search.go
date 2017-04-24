package main

import (
	"fmt"
	"strings"
)

func Search(args []string) {
	// Examples of valid searches:
	// "ingredients=chicken"
	// "tag=soup,vegetarian"

	argString := strings.Join(args[:], " ")

	// Just allow simple searching for now (i.e. only on one field)
	keyValue := strings.Split(argString, "=")

	// strings.Split will always return an array of at least one item
	// (if there are no matches, that item will be an empty string)
	switch len(keyValue) {
	case 2:
		key := keyValue[0]
		value := keyValue[1]
		fmt.Printf("key: %s\nvalue: %s\n", key, value)
	default:
		fmt.Println("Invalid search query.")
	}
}
