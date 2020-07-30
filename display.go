package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func DisplayRecipe(fullFilepath string, displayMetadata bool) {
	recipeFile, err := ParseFile(fullFilepath)
	if err != nil {
		errorMsg := fmt.Sprintf("Unable to read file: %s\n", fullFilepath)
		log.Fatal(errorMsg)
	}
	if displayMetadata == true {
		fmt.Println(string(recipeFile.FrontMatter))
	}
	RenderMarkdown(recipeFile.Markdown)
}
