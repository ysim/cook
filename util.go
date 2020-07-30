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
	key						string
	searchArgs		map[string]Constraint
}

type fmwalk struct {
	fm						map[string][]string
	fullFilepath	string
}

func (w walk) WalkFrontMatter(f func(fmwalk), params ...interface{}) error {
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

			fmwalkData := fmwalk{fm: frontMatter, fullFilepath: fullFilepath}
			f(fmwalkData)
			return nil
		}
	}
	return filepath.Walk(w.prefix, walkGenerator(params))
}

func (w walk) ListKeys(data fmwalk) {
	for k, v := range data.fm {
		key := fmt.Sprintf("%s (%T)", k, v)
		if StringInSlice(key, *w.abstractArray) == false {
			*w.abstractArray = append(*w.abstractArray, key)
		}
	}
}

func (w walk) ListValuesForKey(data fmwalk) {
	valueOfKey := data.fm[w.key]
	for _, v := range valueOfKey {
		if StringInSlice(v, *w.abstractArray) == false {
			*w.abstractArray = append(*w.abstractArray, v)
		}
	}
}

func (w walk) SearchWithArgs(data fmwalk) {
	isMatch := Match(w.searchArgs, data.fm)
	if isMatch {
		fmt.Println(GetBasenameWithoutExt(data.fullFilepath))
	}
}
