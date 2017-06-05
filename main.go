package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"
)

var (
	testArgs []string // Set this to test command line args
	homeDir  string
	prefix   string
	suffix   string
)

type RecipeFile struct {
	FrontMatter []byte
	Markdown    []byte
}

func PrintUsageString() {
	s := `Usage:
    cook [recipe]
    cook edit [recipe]
    cook search key:value[[,|+]value...]
`
	fmt.Printf(s)
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
		if v == nil {
			errMsg := fmt.Sprintf("Key '%s' has the value <nil>", k)
			return nil, errors.New(errMsg)
		}
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
			// Get a new slice with the empty strings removed
			fm[k] = CleanFields(vArray)
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
	var errorMsg string

	fileBytes, err := ioutil.ReadFile(fullFilepath)
	if err != nil {
		return recipeFile, err
	}

	splitBytesArray := bytes.Split(fileBytes, []byte("---"))

	switch {
	case len(splitBytesArray) > 2:
		// If formatted correctly, splitBytesArray[0] is likely an empty string
		assumedFrontMatter := splitBytesArray[1]

		// TODO: This should join splitBytesArray[2:] as Markdown allows for --- to
		// denote a horizontal rule
		assumedMarkdown := splitBytesArray[2]

		// Even if there is only YAML front matter defined with no Markdown
		// content, len(assumedMarkdown) will still be 1 due to a newline
		if len(assumedMarkdown) < 2 {
			errorMsg = "No Markdown has been defined in this file."
			return recipeFile, fmt.Errorf(errorMsg)
		}
		recipeFile := RecipeFile{
			FrontMatter: bytes.TrimSpace(assumedFrontMatter),
			Markdown:    bytes.TrimSpace(assumedMarkdown),
		}
		return recipeFile, nil
	default:
		errorMsg = "Recipe files must consist of a YAML front matter block and non-blank Markdown."
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
		log.WithFields(log.Fields{
			"file": fullFilepath,
		}).Panic(err.Error())
	}
	RenderMarkdown(recipeFile.Markdown)
}

func EditRecipe(fullFilepath string) {
	cmd := exec.Command("vim", fullFilepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
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

	// Log levels
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	var args []string
	switch testArgs {
	case nil:
		flag.Parse()
		args = flag.Args()
	default:
		args = testArgs
	}

	switch len(args) {
	case 0:
		PrintUsageString()
	case 1:
		switch args[0] {
		case "search":
			fmt.Println("Usage: cook search \"key=value\"")
		case "validate":
			ValidateFiles()
		default:
			DisplayRecipe(GetFullFilepath(args[0]))
		}
	default:
		switch args[0] {
		case "search":
			Search(args[1:])
		case "edit":
			EditRecipe(GetFullFilepath(args[1]))
		default:
			PrintUsageString()
		}
	}
}
