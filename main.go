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
	"path"
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

func DisplayRecipe(basename string) {
	prefix := os.Getenv("COOK_RECIPES_DIR")
	homeDir := os.Getenv("HOME")
	if prefix == "" {
		prefix = fmt.Sprintf("%s/.recipes", homeDir)
	}
	suffix := ".md"

	basenameWithSuffix := fmt.Sprintf("%s%s", basename, suffix)
	fullFilepath := path.Join(prefix, basenameWithSuffix)

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

func main() {
	flag.Parse()
	args := flag.Args()

	switch len(args) {
	case 0:
		fmt.Println("Usage:\n\tcook [recipe]\n\tcook search key=value")
	case 1:
		DisplayRecipe(args[0])
	default:
		// Searching front matter.
		switch args[0] {
		case "search":
			fmt.Println("TODO: implement search of front matter")
			//SearchFrontMatter(args[1:])
		default:
			fmt.Printf("No such search term: '%s'\n", args[0])
		}
	}
}
