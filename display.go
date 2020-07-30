package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func DisplayRecipe(fullFilepath string) {
	recipeFile, err := ParseFile(fullFilepath)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable to read file: %s\n", fullFilepath)
		log.Fatal(errorMsg)
	}
	RenderMarkdown(recipeFile.Markdown)
}
