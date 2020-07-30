package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
	binary   string
)

type RecipeFile struct {
	FrontMatter []byte
	Markdown    []byte
}

func PrintUsageString() {
	s := `Usage:
    cook [recipe]
    cook edit [recipe]
    cook list [-key=KEY]
    cook new -filename=FILENAME
    cook search key:value
    cook validate [recipe]
    cook version
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
		errorMsg := fmt.Sprintf("Unable to read file: %s\n", fullFilepath)
		log.Fatal(errorMsg)
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

	binary = os.Args[0]
	switch binary {
	case "cook":
		prefix = os.Getenv("COOK_RECIPES_DIR")
		if prefix == "" {
			prefix = fmt.Sprintf("%s/.recipes", homeDir)
		}
	case "concoct":
		prefix = os.Getenv("CONCOCT_RECIPES_DIR")
		if prefix == "" {
			prefix = fmt.Sprintf("%s/.drinks", homeDir)
		}
	default:
		// Doesn't matter what the prefix is if the binary is neither `cook`
		// nor `concoct` as long as the the path exists; this likely means
		// that we're in test mode.
		prefix = homeDir
	}

	fileInfo, err := os.Lstat(prefix)
	if err != nil {
		errorMsg := fmt.Sprintf("There was an error getting file info for the prefix: %s\n", prefix)
		log.Fatal(errorMsg)
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		prefix, err = os.Readlink(prefix)
		if err != nil {
			errorMsg := fmt.Sprintf("Unable to read symlink: %s\n", prefix)
			log.Fatal(errorMsg)
		}
	}
	suffix = os.Getenv("COOK_RECIPES_EXT")
	if suffix == "" {
		suffix = ".md"
	}

	// Log levels
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.FatalLevel)
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
	newSubcommandFlagset := flag.NewFlagSet("newFlagset", flag.ContinueOnError)
	listSubcommandFlagset := flag.NewFlagSet("listFlagset", flag.ContinueOnError)

	// `new` subcommand flag pointers
	recipeFilename := newSubcommandFlagset.String("filename", "", "The recipe filename (without the extension).")
	recipeName := newSubcommandFlagset.String("name", "", "The recipe name.")
	newSubcommandFlagset.Var(&fieldFlags, "f", "An arbitrary field.")

	// `list` subcommand flag pointers
	keyPtr := listSubcommandFlagset.String("key", "", "The key name for which to list values.")

	args := os.Args
	switch len(args) {
	case 1:
		PrintUsageString()
	case 2:
		switch args[1] {
		case "edit", "new", "search":
			PrintUsageString()
		case "list":
			List("")
		case "validate":
			ValidateFiles()
		case "version":
			PrintVersion()
		default:
			DisplayRecipe(GetFullFilepath(args[1]))
		}
	default:
		switch args[1] {
		case "edit":
			EditRecipe(GetFullFilepath(args[2]))
		case "new":
			newSubcommandFlagset.Parse(args[2:])
			CreateNewRecipe(*recipeFilename, *recipeName, fieldFlags)
		case "list":
			listSubcommandFlagset.Parse(args[2:])
			List(*keyPtr)
		case "search":
			Search(args[2:])
		case "validate":
			ValidateSingleFile(GetFullFilepath(args[2]))
		default:
			PrintUsageString()
		}
	}
}
