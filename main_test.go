package main

import (
	"github.com/stretchr/testify/assert"
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
	}

	for _, table := range tables {
		actualResult, actualError := ParseFrontMatter(table.FrontMatterBytes)
		assert.Equal(t, table.ExpectedResult, actualResult)
		assert.Equal(t, table.ExpectedError, actualError)
	}
}
