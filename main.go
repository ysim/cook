package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type FrontMatter struct {
	Name        string   `yaml:"name"`
	Tags        []string `yaml:"tags"`
	Ingredients []string `yaml:"ingredients"`
}

func ParseFrontMatter(fmBytes []byte) {
	fm := FrontMatter{}
	err := yaml.Unmarshal([]byte(fmBytes), &fm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", fm)
}

func ParseMarkdown(mdBytes []byte) {
	output := blackfriday.MarkdownBasic(mdBytes)
	fmt.Println(string(output))
}

func main() {
	prefix := os.Getenv("COOK_RECIPES_DIR")
	if prefix == "" {
		prefix = "${HOME}/.recipes"
	}
	suffix := ".md"

	flag.Parse()
	file := flag.Args()[0]
	fullFilepath := fmt.Sprintf("%s/%s%s", prefix, file, suffix)

	fileBytes, err := ioutil.ReadFile(fullFilepath)
	if err != nil {
		panic(err)
	}

	splitBytesArray := bytes.Split(fileBytes, []byte("---"))
	markdownBytes := splitBytesArray[len(splitBytesArray)-1]

	// No YAML front matter has been defined
	if len(splitBytesArray) < 2 {
		frontMatterBytes := splitBytesArray[1]
		ParseFrontMatter(frontMatterBytes)
	}

	ParseMarkdown(markdownBytes)
}
