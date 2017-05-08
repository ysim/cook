package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
