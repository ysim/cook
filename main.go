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

var (
	homeDir string
	prefix  string
	suffix  string
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
	// TODO: render html in the command line
	fmt.Println(string(output))
}

func ParseFile(basename string) [][]byte {
	basenameWithSuffix := fmt.Sprintf("%s%s", basename, suffix)
	fullFilepath := path.Join(prefix, basenameWithSuffix)

	fileBytes, err := ioutil.ReadFile(fullFilepath)
	if err != nil {
		panic(err)
	}

	splitBytesArray := bytes.Split(fileBytes, []byte("---"))
	return splitBytesArray
}

func DisplayRecipe(basename string) {
	splitBytesArray := ParseFile(basename)
	markdownBytes := splitBytesArray[len(splitBytesArray)-1]

	// No YAML front matter has been defined
	if len(splitBytesArray) < 2 {
		frontMatterBytes := splitBytesArray[1]
		ParseFrontMatter(frontMatterBytes)
	}

	ParseMarkdown(markdownBytes)
}

func main() {
	homeDir = os.Getenv("HOME")
	prefix = os.Getenv("COOK_RECIPES_DIR")
	if prefix == "" {
		prefix = fmt.Sprintf("%s/.recipes", homeDir)
	}
	suffix = ".md"

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
