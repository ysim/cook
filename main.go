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

func PrintUsageString() {
	fmt.Printf("Usage:\n\tcook [recipe]\n\tcook search key=value\n")
}

func ParseFrontMatter(fmBytes []byte) {
	var fm interface{}
	err := yaml.Unmarshal([]byte(fmBytes), &fm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", fm)
}

func RenderMarkdown(mdBytes []byte) {
	// TODO: decide how to render the markdown
	output := blackfriday.MarkdownBasic(mdBytes)
	fmt.Println(string(output))
}

// TODO: Modify to accept the full filepath as the only argument
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
	// TODO: How should we handle instances where the user has created a
	// Markdown file that doesn't conform to the standard format?
	splitBytesArray := ParseFile(basename)
	markdownBytes := splitBytesArray[len(splitBytesArray)-1]
	RenderMarkdown(markdownBytes)
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
		PrintUsageString()
	case 1:
		switch args[0] {
		case "search":
			fmt.Println("Usage: cook search \"key=value\"")
		default:
			DisplayRecipe(args[0])
		}
	default:
		// Searching front matter.
		switch args[0] {
		case "search":
			Search(args[1:])
		default:
			PrintUsageString()
		}
	}
}
