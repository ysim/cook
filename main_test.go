package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	tables := []struct {
		FrontMatterBytes []byte
		ExpectedResult   map[string][]string
		ExpectedError    error
	}{
		{
			[]byte("name: pasta salad\ntags: [summer, salad, pasta]"),
			map[string][]string{"name": []string{"pasta salad"}, "tags": []string{"summer", "salad", "pasta"}},
			nil,
		},
		{
			[]byte("name: nachos\ningredients: [minced beef]"),
			map[string][]string{"name": []string{"nachos"}, "ingredients": []string{"minced beef"}},
			nil,
		},
		{
			[]byte("name: \ningredients: [a mystery]"),
			nil,
			fmt.Errorf("Key 'name' has the value <nil>"),
		},
	}

	for _, table := range tables {
		actualResult, actualError := ParseFrontMatter(table.FrontMatterBytes)
		assert.Equal(t, table.ExpectedResult, actualResult)
		assert.Equal(t, table.ExpectedError, actualError)
	}
}

func CreateTempTestFile(contents []byte) *os.File {
	// Specifying an empty string for the first arg means that Tempfile will
	// use the default directory for temporary files
	tmpfile, err := ioutil.TempFile("", "testrecipe")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(contents); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpfile
}

func TestParseFile(t *testing.T) {
	tables := []struct {
		FileContentString []byte
		ExpectedResult    RecipeFile
		ExpectedError     error
	}{
		{
			[]byte("---\nname: french fries\ntexture: [crispy]\n---\nThe recipe."),
			RecipeFile{[]byte("name: french fries\ntexture: [crispy]"), []byte("The recipe.")},
			nil,
		},
		{
			[]byte("---\nname: lasagna\n---\n---\nThere was no Markdown."),
			RecipeFile{},
			fmt.Errorf("No Markdown has been defined in this file."),
		},
	}

	for _, table := range tables {
		f := CreateTempTestFile(table.FileContentString)
		actualResult, actualError := ParseFile(f.Name())
		assert.Equal(t, table.ExpectedResult, actualResult)
		assert.Equal(t, table.ExpectedError, actualError)
		os.Remove(f.Name()) // Remember to clean up!
	}
}
