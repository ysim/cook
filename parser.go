package main

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"strings"
)

const (
	style_h1    = "\x1b[1;4;92m" // green, bold, underline, high intensity
	style_h2    = "\x1b[4;92m"   // green, underline, high intensity
	style_h3    = "\x1b[1;92m"   // green, bold, high intensity
	style_h4    = "\x1b[1;93m"   // yellow, bold, high intensity
	style_h5    = "\x1b[0;93m"   // yellow, high intensity
	style_h6    = "\x1b[0;33m"   // yellow
	style_reset = "\x1b[0m"
)

func ParseHTML(htmlBytes []byte) ([]string, error) {
	var output []string
	var listPosition int
	var isOrderedList bool

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
				switch tagName {
				case "h1":
					output = append(output, style_h1)
				case "h2":
					output = append(output, style_h2)
				case "h3":
					output = append(output, style_h3)
				case "h4":
					output = append(output, style_h4)
				case "h5":
					output = append(output, style_h5)
				case "h6":
					output = append(output, style_h6)
				case "li":
					listPosition++
					switch isOrderedList {
					case true:
						output = append(output, fmt.Sprintf("%d. ", listPosition))
					case false:
						output = append(output, "- ")
					}
				case "ol":
					isOrderedList = true
					listPosition = 0
					output = append(output, "\n")
				}
			case html.EndTagToken:
				switch tagName {
				case "h1", "h2", "h3", "h4", "h5", "h6":
					output = append(output, style_reset)
				case "ol":
					isOrderedList = false
					listPosition = 0
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
