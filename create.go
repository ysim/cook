package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

const (
	keyValueDelimiter = "="
	openBracket       = "["
	closeBracket      = "]"
	newRecipeTemplate = "---\nname: {{.Name}}{{range .Fields}}\n{{.Key}}: {{.Value}}{{end}}\nsource: [\"\"]\n---\n# {{.Name}}\n\n## INGREDIENTS\n\n## INSTRUCTIONS\n"
)

type RecipeField struct {
	Key   string
	Value string
}

type RecipeSkeleton struct {
	Name   string
	Fields []RecipeField
}

func validateFieldValue(value string) string {
	// First, validate that we were given something in a format that we expect.
	valueContent := strings.Trim(value, strings.Join([]string{openBracket, closeBracket}, ""))
	valueArray := strings.Split(valueContent, ",")
	var trimmedValueArray []string
	for _, s := range valueArray {
		trimmedValueArray = append(trimmedValueArray, strings.TrimSpace(s))
	}

	// Then, put it all back together in a standard format.
	normalizedFieldValue := fmt.Sprintf("[%s]", strings.Join(trimmedValueArray, ", "))

	return normalizedFieldValue
}

func validateFields(fieldFlags []string) (map[string]interface{}, error) {
	fields := make(map[string]interface{})
	for _, rawField := range fieldFlags {
		splitValueArray := strings.Split(rawField, keyValueDelimiter)

		// If there's a third item in the array, we don't care about it. :)
		key := splitValueArray[0]
		value := splitValueArray[1]

		fields[key] = validateFieldValue(value)
	}
	return fields, nil
}

func validateNewRecipe(filename string) string {
	if filename == "" {
		errMsg := "You must provide at least a filename in order to create a new recipe."
		log.Fatal(errMsg)
	}

	if strings.HasSuffix(filename, suffix) {
		filename = strings.TrimSuffix(filename, suffix)
	}

	filepath := GetFullFilepath(filename)
	_, err := os.Stat(filepath)

	// Not the same as os.IsExist! This is because os.Stat doesn't throw an
	// error if the file exists, so os.IsExist would receive a nil value for err.
	// So, an os.IsExist(err) block would never execute for a file that exists.
	if !os.IsNotExist(err) {
		errMsg := fmt.Sprintf("There already exists a file at the path: %s\n", filepath)
		log.Fatal(errMsg)
	}
	return filepath
}

func writeNewRecipeFile(filepath string, name string, fields map[string]interface{}) {
	var recipeFields []RecipeField
	for k, v := range fields {
		recipeFields = append(recipeFields, RecipeField{Key: k, Value: v.(string)})
	}

	recipeVariables := RecipeSkeleton{
		Name:   name,
		Fields: recipeFields,
	}

	f, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("An error occurred creating a file at: %s\n", filepath)
	}
	defer f.Close()

	tpl, err := template.New("newrecipe").Parse(newRecipeTemplate)
	if err != nil {
		log.Fatal("Error parsing recipe template.")
	}
	err = tpl.Execute(f, recipeVariables)
	if err != nil {
		log.Fatal("Error executing template.")
	}
}

func CreateNewRecipe(filename string, name string, fieldFlags []string) {
	filepath := validateNewRecipe(filename)

	var validatedFields map[string]interface{}
	var validateFieldsErr error
	if len(fieldFlags) > 0 {
		validatedFields, validateFieldsErr = validateFields(fieldFlags)
		if validateFieldsErr != nil {
			log.Fatal("An error occurred while validating the new recipe fields.")
		}
	}

	writeNewRecipeFile(filepath, name, validatedFields)
	EditRecipe(filepath)
}
