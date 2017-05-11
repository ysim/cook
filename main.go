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

type RecipeFile struct {
	FrontMatter []byte
	Markdown    []byte
}

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

	// Now make a type assertion into map[string][]string to make querying
	// easier
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

func ParseFile(fullFilepath string) (RecipeFile, error) {
	var recipeFile RecipeFile
	errorMsg := fmt.Sprintf("This file could not be parsed: %s", fullFilepath)
	fileBytes, err := ioutil.ReadFile(fullFilepath)
	if err != nil {
		return recipeFile, err
	}

	splitBytesArray := bytes.Split(fileBytes, []byte("---"))

	switch {
	case len(splitBytesArray) > 2:
		// If formatted correctly, splitBytesArray[0] is likely an empty string
		assumedFrontMatter := splitBytesArray[1]
		assumedMarkdown := splitBytesArray[2]

		// Even if there is only YAML front matter defined with no Markdown
		// content, len(assumedMarkdown) will still be 1 due to a newline
		if len(assumedMarkdown) < 2 {
			return recipeFile, fmt.Errorf(errorMsg)
		}
		recipeFile := RecipeFile{
			FrontMatter: assumedFrontMatter,
			Markdown:    assumedMarkdown,
		}
		return recipeFile, nil
	default:
		return recipeFile, fmt.Errorf(errorMsg)
	}
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
	recipeFile, err := ParseFile(fullFilepath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	RenderMarkdown(recipeFile.Markdown)
}

func init() {
	homeDir = os.Getenv("HOME")
	prefix = os.Getenv("COOK_RECIPES_DIR")
	if prefix == "" {
		prefix = fmt.Sprintf("%s/.recipes", homeDir)
	}
	suffix = os.Getenv("COOK_RECIPES_EXT")
	if suffix == "" {
		suffix = ".md"
	}
}

func main() {
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
