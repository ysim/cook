package main

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"strings"
)

const (
	style_h1     = "\n\n\x1b[1;4;92m" // green, bold, underline, high intensity
	style_h2     = "\n\n\x1b[4;92m"   // green, underline, high intensity
	style_h3     = "\n\n\x1b[1;92m"   // green, bold, high intensity
	style_h4     = "\n\n\x1b[1;93m"   // yellow, bold, high intensity
	style_h5     = "\n\n\x1b[0;93m"   // yellow, high intensity
	style_h6     = "\n\n\x1b[0;33m"   // yellow
	style_strong = "\x1b[1;4;37m"     // grey, bold, underline
	style_em     = "\x1b[4;37m"       // grey, underline
	style_reset  = "\x1b[0m"
)

func ParseHTML(htmlBytes []byte) ([]string, error) {
	var output []string
	var listPosition int
	var isOrderedList bool
	var inListItem bool
	var depth int

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
				case "strong":
					output = append(output, style_strong)
				case "em":
					output = append(output, style_em)
				case "p":
					if !inListItem {
						output = append(output, "\n\n")
					}
				case "li":
					inListItem = true
					listPosition++
					output = append(output, "\n")
					if depth > 1 {
						output = append(output, strings.Repeat(" ", depth*2))
					}
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
					depth++
				case "ul":
					depth++
				}
			case html.EndTagToken:
				switch tagName {
				case "h1", "h2", "h3", "h4", "h5", "h6", "strong", "em":
					output = append(output, style_reset)
				case "ol":
					isOrderedList = false
					listPosition = 0
					depth--
				case "ul":
					depth--
				case "li":
					inListItem = false
				}
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
