package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func StringInSlice(s string, slice []string) bool {
	for _, i := range slice {
		if s == i {
			return true
		}
	}
	return false
}

func PrintArrayOnNewlines(a []string) {
	for _, v := range a {
		fmt.Println(v)
	}
}

type walk struct {
	prefix				string
	abstractArray	*[]string
}

func (w walk) WalkFrontMatter(f func(map[string][]string) error, params ...interface{}) error {
	walkGenerator := func(...interface{}) filepath.WalkFunc {
		return func(fullFilepath string, info os.FileInfo, err error) error {
			shouldSkip := ShouldSkipFile(info, err)
			if shouldSkip {
				return nil
			}

			recipeFile, err := ParseFile(fullFilepath)
			if err != nil {
				log.WithFields(log.Fields{
					"file": fullFilepath,
				}).Warn(err.Error())
				return nil
			}

			frontMatter, err := ParseFrontMatter(recipeFile.FrontMatter)
			if err != nil {
				log.WithFields(log.Fields{
					"file": fullFilepath,
				}).Warn("Unknown type detected in front matter")
			}

			f(frontMatter)
			return nil
		}
	}
	return filepath.Walk(w.prefix, walkGenerator(params))
}

func (w walk) WalkListKeys(frontMatter map[string][]string) error {
	for k, v := range frontMatter {
		key := fmt.Sprintf("%s (%T)", k, v)
		if StringInSlice(key, *w.abstractArray) == false {
			*w.abstractArray = append(*w.abstractArray, key)
		}
	}
	return nil
}
