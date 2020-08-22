package testhelpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// TmplSuffix is the file suffix for template files
var TmplSuffix = ".tmpl"

// RegSuffix is the file suffix for non-template files
var RegSuffix = ".txt"

// CreateTestTemplateProject -
func CreateTestTemplateProject(projectDir, subDir, outputDir string) (string, string, string, string, error) {
	subDirPath := filepath.Join(projectDir, subDir)
	tmplFile := "template" + TmplSuffix
	regFile := "file" + RegSuffix

	if err := os.MkdirAll(subDirPath, 0755); err != nil {
		return "", "", "", "", err
	}
	if err := CreateFileWithContents(filepath.Join(subDirPath, tmplFile), "{{/* a comment */}}"); err != nil {
		return "", "", "", "", err
	}
	if err := CreateFileWithContents(filepath.Join(subDirPath, regFile), "notemplating"); err != nil {
		return "", "", "", "", err
	}
	return filepath.Join(subDirPath, tmplFile), filepath.Join(subDirPath, regFile), filepath.Join(outputDir, subDir, tmplFile), filepath.Join(outputDir, subDir, regFile), nil
}

// CreateFileWithContents -
func CreateFileWithContents(path, contents string) error {
	return ioutil.WriteFile(path, []byte(contents), 0644)
}
