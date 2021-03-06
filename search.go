package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

const (
	q_kv_separator = ":"
	q_value_or     = ","
	q_value_and    = "+"
)

// For now, we will only allow positive matches (i.e. for matches and not
// non-matches).
type Constraint struct {
	Terms        []string
	Relationship string // and, or
}

func CleanFields(a []string) []string {
	var cleanedSlice []string
	for _, element := range a {
		element = strings.TrimSpace(element)
		if len(element) > 0 {
			cleanedSlice = append(cleanedSlice, element)
		}
	}
	return cleanedSlice
}

func ProcessValueAnd(queryTerms []string, fileTerms []string) bool {
	for _, queryTerm := range queryTerms {
		match := false
		for _, fileTerm := range fileTerms {
			if strings.Contains(fileTerm, queryTerm) {
				match = true
			}
		}
		if !match {
			return false
		}
	}
	return true
}

// TODO: This could be made more efficient
func ProcessValueOr(queryTerms []string, fileTerms []string) bool {
	for _, queryTerm := range queryTerms {
		for _, fileTerm := range fileTerms {
			if strings.Contains(fileTerm, queryTerm) {
				return true
			}
		}
	}
	return false
}

// Determine whether there is a match between a file and the search query.
// If the relationship is 'or', exit the moment a match is found.
// If the relationship is 'and', exit when all the conditions are met.
func Match(queryArgs map[string]Constraint, fileArgs map[string][]string) bool {
	for queryKey, queryValue := range queryArgs {
		// ok is set to true if the key exists, false if not
		fileValueArray, ok := fileArgs[queryKey]

		if ok {
			switch queryValue.Relationship {
			case "or":
				return ProcessValueOr(queryValue.Terms, fileValueArray)
			case "and":
				return ProcessValueAnd(queryValue.Terms, fileValueArray)
			}
		}
	}
	return false
}

// Implements some common rules for filepath.Walkfunc for which kinds of
// files will be skipped.
func ShouldSkipFile(info os.FileInfo, err error) bool {
	if err != nil {
		log.WithFields(log.Fields{
			"file": info.Name(),
		}).Warn(err)
		return true
	}
	// Don't descend into directories for now
	if info.IsDir() {
		return true
	}
	// Ignore hidden files
	if strings.HasPrefix(info.Name(), ".") {
		return true
	}
	return false
}

func GetFieldValueLogic(s string) (Constraint, error) {
	var cleanedField Constraint
	// Make sure there are only , or &
	andCount := strings.Count(s, q_value_and)
	orCount := strings.Count(s, q_value_or)
	switch {
	case andCount == 0 && orCount == 0:
		cleanedString := strings.TrimSpace(s)
		cleanedField = Constraint{[]string{cleanedString}, "or"}
	case andCount > 0 && orCount == 0:
		fields := CleanFields(strings.Split(s, q_value_and))
		cleanedField = Constraint{fields, "and"}
	case andCount == 0 && orCount > 0:
		fields := CleanFields(strings.Split(s, q_value_or))
		cleanedField = Constraint{fields, "or"}
	default:
		errMsg := fmt.Sprintf("Only one type of '%s' and '%s' allowed in a search query", q_value_and, q_value_or)
		return cleanedField, errors.New(errMsg)
	}
	return cleanedField, nil
}

func ParseSearchQuery(args []string) (map[string]Constraint, error) {
	q := make(map[string]Constraint)

	// Consolidate all arguments
	argString := strings.Join(args[:], " ")

	// TODO: Replace with custom function once multifield search is figured out
	// Now split into fields
	fields := strings.Fields(argString)

	// TODO: Remove the slice [:1] once multifield search is figured out
	for _, f := range fields[:1] {
		// strings.Split will always return an array of at least one item
		// (if there are no matches, that item will be an empty string)
		splitField := strings.Split(f, q_kv_separator)
		if len(splitField) != 2 {
			errMsg := fmt.Sprintf("Exactly one '%s' is required per whitespace-delimited argument", q_kv_separator)
			return nil, errors.New(errMsg)
		}

		key, value := splitField[0], splitField[1]
		valueLogic, err := GetFieldValueLogic(value)
		if err != nil {
			return nil, err
		}
		q[key] = valueLogic
	}
	return q, nil
}

func Search(args []string) {
	parsedQuery, parseErr := ParseSearchQuery(args)
	if parseErr != nil {
		errMsg := fmt.Sprintf("Invalid search query: %s\n", parseErr.Error())
		log.Fatal(errMsg)
	}

	w := walk{prefix: prefix, searchArgs: parsedQuery}
	searchErr := w.WalkFrontMatter(w.SearchWithArgs)
	if searchErr != nil {
		log.Warn(searchErr)
	}
}
