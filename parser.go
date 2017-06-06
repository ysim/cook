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
			output = append(output, strings.TrimSpace(string(z.Text())))
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			tagName := string(tn)
			switch token {
			case html.StartTagToken:
				if tagName == "h1" {
					output = append(output, "\x1b[1;4;92m")
				}
			case html.EndTagToken:
				if tagName == "h1" {
					output = append(output, "\x1b[0m")
				}
				output = append(output, "\n")
			}
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