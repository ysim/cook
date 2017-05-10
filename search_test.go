package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatch(t *testing.T) {
	tables := []struct {
		QueryArgs      map[string]Constraint
		FileArgs       map[string][]string
		ExpectedResult bool
	}{
		{
			map[string]Constraint{
				"tags": Constraint{[]string{"vegan", "summer"}, "or"},
			},
			map[string][]string{"tags": []string{"vegan", "soup"}},
			true,
		},
		{
			map[string]Constraint{
				"ingredients": Constraint{[]string{"salmon"}, "or"},
			},
			map[string][]string{"ingredients": []string{"carrots", "salmon"}},
			true,
		},
	}

	for _, table := range tables {
		actualResult := Match(table.QueryArgs, table.FileArgs)
		verboseDescription := fmt.Sprintf("QueryArgs: %s\nFileArgs: %s\n", table.QueryArgs, table.FileArgs)
		assert.Equal(t, table.ExpectedResult, actualResult, verboseDescription)
	}
}

func TestCleanFields(t *testing.T) {
	tables := []struct {
		Arg            []string
		ExpectedResult []string
	}{
		{
			[]string{"zucchini", "eggplant", "peppers"},
			[]string{"zucchini", "eggplant", "peppers"},
		},
		{
			[]string{"minced beef ", "  carrots", " potatoes "},
			[]string{"minced beef", "carrots", "potatoes"},
		},
		{
			[]string{"", " ", "		", "  ", "avocadoes"},
			[]string{"avocadoes"},
		},
	}

	for _, table := range tables {
		actualResult := CleanFields(table.Arg)
		assert.Equal(t, table.ExpectedResult, actualResult)
	}
}

func TestParseSearchQuery(t *testing.T) {
	tables := []struct {
		Args           []string
		ExpectedResult map[string]Constraint
		ExpectedError  error
	}{
		{
			[]string{"tags:breakfast"},
			map[string]Constraint{
				"tags": Constraint{
					[]string{"breakfast"},
					"or",
				},
			},
			nil,
		},
		{
			[]string{"tags:soup,vegetarian"},
			map[string]Constraint{
				"tags": Constraint{
					[]string{"soup", "vegetarian"},
					"or",
				},
			},
			nil,
		},
	}

	for _, table := range tables {
		actualResult, _ := ParseSearchQuery(table.Args)
		assert.Equal(t, table.ExpectedResult, actualResult)
	}
}
