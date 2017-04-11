package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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
	flag.Parse()
	file := flag.Args()[0]

	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	splitBytesArray := bytes.Split(fileBytes, []byte("---"))
	frontMatterBytes := splitBytesArray[1]
	markdownBytes := splitBytesArray[2]

	ParseFrontMatter(frontMatterBytes)
	ParseMarkdown(markdownBytes)
}
