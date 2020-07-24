package main

import (
	"fmt"
)

func PrintArrayOnNewlines(a []string) {
	for _, v := range a {
		fmt.Println(v)
	}
}
