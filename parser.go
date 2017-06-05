package main

import (
	"fmt"
	"github.com/russross/blackfriday"
)

func RenderMarkdown(mdBytes []byte) error {
	// TODO: decide how to render the markdown
	output := blackfriday.MarkdownBasic(mdBytes)
	fmt.Println(output)
	return nil
}
