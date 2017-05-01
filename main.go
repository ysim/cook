package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
)

var (
	homeDir string
	prefix  string
	suffix  string
)

func PrintUsageString() {
	fmt.Printf("Usage:\n\tcook [recipe]\n\tcook search key=value\n")
}

func ParseFrontMatter(fmBytes []byte) (map[string][]string, error) {
	// Unmarshal into ...interface{} initially to allow for flexible data
	// structures
	var rfm map[string]interface{}
	err := yaml.Unmarshal([]byte(fmBytes), &rfm)
	if err != nil {
		return nil, err
	}

	fm := make(map[string][]string, len(rfm))
	for k, v := range rfm {
		t := reflect.TypeOf(v).Kind()
		switch t {
		case reflect.String:
			fm[k] = []string{v.(string)}
		case reflect.Slice:
			coercedArray := v.([]interface{})
			vArray := make([]string, len(coercedArray))
			for _, item := range coercedArray {
				vArray = append(vArray, item.(string))
			}
			fm[k] = vArray
		default:
			return nil, errors.New("Type was not string or slice.")
		}
	}
	return fm, nil
}

func RenderMarkdown(mdBytes []byte) {
	// TODO: decide how to render the markdown
	output := blackfriday.MarkdownBasic(mdBytes)
	fmt.Println(string(output))
}

func ParseFile(fullFilepath string) [][]byte {
	fileBytes, err := ioutil.ReadFile(fullFilepath)
	if err != nil {
		panic(err)
	}

	splitBytesArray := bytes.Split(fileBytes, []byte("---"))
	return splitBytesArray
}

func GetFullFilepath(basename string) string {
	basenameWithSuffix := fmt.Sprintf("%s%s", basename, suffix)
	fullFilepath := path.Join(prefix, basenameWithSuffix)
	return fullFilepath
}

func GetBasenameWithoutExt(fullFilepath string) string {
	return strings.Replace(
		path.Base(fullFilepath),
		path.Ext(fullFilepath),
		"",
		-1)
}

func DisplayRecipe(fullFilepath string) {
	// TODO: How should we handle instances where the user has created a
	// Markdown file that doesn't conform to the standard format?
	splitBytesArray := ParseFile(fullFilepath)
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
			recipeBasenameWithoutExtension := args[0]
			recipeFullFilepath := GetFullFilepath(recipeBasenameWithoutExtension)
			DisplayRecipe(recipeFullFilepath)
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
