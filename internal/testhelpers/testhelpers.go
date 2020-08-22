package testhelpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// CreateTestTemplateProject -
func CreateTestTemplateProject(projectDir, subDir, outputDir string) (string, string, string, string, error) {
	subDirPath := filepath.Join(projectDir, subDir)
	tmplFile := "template.txt"
	regFile := "file.txt"

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
