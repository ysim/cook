package main

import (
	"fmt"
)

func StringInSlice(s string, slice []string) bool {
	for _, i := range slice {
		if s == i {
			return true
		}
	}
	return false
}

func PrintArrayOnNewlines(a []string) {
	for _, v := range a {
		fmt.Println(v)
	}
}
