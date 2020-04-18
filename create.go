package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	keyValueDelimiter = "="
	openBracket = "["
	closeBracket = "]"
)

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
	// TODO: Be able to use a delimiter other than =, like :
	// TODO: Be able to use brackets other than []

	fields := make(map[string]interface{})
	// How do we determine if something is a list? Square brackets? A comma at the end?
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
		fmt.Println("You must provide at least a filename in order to create a new recipe.")
		os.Exit(1)
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
		fmt.Printf("There already exists a file at the path: %s\n", filepath)
		os.Exit(1)
	}
	return filepath
}

func writeNewRecipeFile(filepath string, name string, fields map[string]interface{}) {
	fmt.Printf("New recipe file at path: %s\n", filepath)
	fmt.Printf("Recipe name: %s\n", name)
	fmt.Printf("Recipe fields: %q\n", fields)
}

func CreateNewRecipe(filename string, name string, fieldFlags []string) {
	filepath := validateNewRecipe(filename)
	validatedFields, validateFieldsErr := validateFields(fieldFlags)
	if validateFieldsErr != nil {
		fmt.Println("An error occurred while validating the new recipe fields.")
		os.Exit(1)
	}
	writeNewRecipeFile(filepath, name, validatedFields)
}
