package main

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"strings"
)

func ParseHTML(htmlBytes []byte) ([]string, error) {
	var output []string
	z := html.NewTokenizer(bytes.NewReader(htmlBytes))
	for {
		token := z.Next()
		switch token {
		case html.ErrorToken:
			return output, z.Err()
		case html.TextToken:
			s := z.Text()
			output = append(output, string(s))
		}
	}
}

func RenderMarkdown(mdBytes []byte) error {
	htmlBytes := blackfriday.MarkdownBasic(mdBytes)
	outputSlice, err := ParseHTML(htmlBytes)
	if err.Error() != "EOF" {
		return err
	}
	fmt.Println(strings.Join(outputSlice, ""))
	return nil
}
