package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		ExpectedResult map[string][]string
		ExpectedError  error
	}{
		{
			[]string{"tags:breakfast"},
			map[string][]string{"tags": []string{"breakfast"}},
			nil,
		},
		{
			[]string{"tags:soup,vegetarian"},
			map[string][]string{"tags": []string{"soup", "vegetarian"}},
			nil,
		},
		{
			[]string{"tags:soup,vegetarian ingredients:cauliflower"},
			map[string][]string{
				"tags":        []string{"soup", "vegetarian"},
				"ingredients": []string{"cauliflower"},
			},
			nil,
		},
	}

	for _, table := range tables {
		actualResult, _ := ParseSearchQuery(table.Args)
		assert.Equal(t, table.ExpectedResult, actualResult)
	}
}
