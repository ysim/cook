package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
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
	version  string
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

func PrintVersion() {
	fmt.Println(version)
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
		errorMsg = "Recipe files must consist of a YAML front matter block delimited by three hyphens (---) followed non-blank Markdown."
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
		fmt.Printf("Unable to read file: %s\n", fullFilepath)
		os.Exit(1)
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
	fileInfo, err := os.Lstat(prefix)
	if err != nil {
		fmt.Printf("There was an error getting file info for: %s\n", prefix)
		os.Exit(1)
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		prefix, err = os.Readlink(prefix)
		if err != nil {
			fmt.Printf("Unable to read symlink: %s\n", prefix)
			os.Exit(1)
		}
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

type flagArray []string

func (i *flagArray) String() string {
  return "Spooky string."
}

func (i *flagArray) Set(value string) error {
	// While it is possible (and more efficient) to create a map while we're
	// iterating through the flags, I'd prefer to keep the parser logic separate
	// from the code here, which is just concerned with CLI processing.
  *i = append(*i, value)
  return nil
}

var fieldFlags flagArray

func main() {
	recipeFilename := flag.String("filename", "new-recipe", "The recipe filename (without the extension).")
	recipeName := flag.String("name", "New Recipe", "The recipe name.")
	flag.Var(&fieldFlags, "f", "An arbitrary field.")

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
		case "new":
			CreateNewRecipe(*recipeFilename, *recipeName, fieldFlags)
		case "search":
			fmt.Println("Usage: cook search \"key:value\"")
		case "validate":
			ValidateFiles()
		case "version":
			PrintVersion()
		default:
			DisplayRecipe(GetFullFilepath(args[0]))
		}
	default:
		switch args[0] {
		case "search":
			Search(args[1:])
		case "edit":
			EditRecipe(GetFullFilepath(args[1]))
		case "validate":
			ValidateSingleFile(GetFullFilepath(args[1]))
		default:
			PrintUsageString()
		}
	}
}
