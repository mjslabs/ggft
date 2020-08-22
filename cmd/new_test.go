package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mjslabs/ggft/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestCmdNew(t *testing.T) {
	t.Run("NewProject", testNewProject)
	t.Run("NewProjectNotFound", testNewProjectNotFound)
	t.Run("ProcessPath", testProcessPath)
	t.Run("CopyFile", testCopyFile)
	t.Run("ShouldCopyFile", testShouldCopyFile)
	t.Run("ShouldCopyAndTrimFile", testShouldTrimFile)
	t.Run("ShouldCopyAndTrimFile", testShouldIgnore)
}

func testNewProject(t *testing.T) {
	projectDir := "testingInPath"
	subDir := "testingSubDir"
	outputDir := "testingOutPath"

	inTmplFile, inRegFile, outTmplFile, outRegFile, err := testhelpers.CreateTestTemplateProject(projectDir, subDir, outputDir)
	assert.NoError(t, err)
	assert.NoError(t, newProject(projectDir, outputDir, []string{}, []string{}, []string{}, []string{}))

	hashOne, err := md5OfFile(inTmplFile)
	assert.NoError(t, err)
	hashTwo, err := md5OfFile(outTmplFile)
	assert.NoError(t, err)
	assert.NotEqual(t, hashOne, hashTwo)

	// Check that the template changed
	b, e := filesAreTheSame(inTmplFile, outTmplFile)
	assert.NoError(t, e)
	assert.False(t, b)

	// Check that the non-template is the same
	b, e = filesAreTheSame(inRegFile, outRegFile)
	assert.NoError(t, e)
	assert.True(t, b)

	assert.NoError(t, os.RemoveAll(projectDir))
	assert.NoError(t, os.RemoveAll(outputDir))

	// test for -s .tmpl
	inTmplFile, inRegFile, outTmplFile, outRegFile, err = testhelpers.CreateTestTemplateProject(projectDir, subDir, outputDir)
	assert.NoError(t, newProject(projectDir, outputDir, []string{testhelpers.TmplSuffix}, []string{}, []string{}, []string{}))
	// Check that the file exists
	assert.FileExists(t, outTmplFile)
	// Check that outTmplFile didn't get parsed as a template
	b, e = filesAreTheSame(inTmplFile, outTmplFile)
	assert.NoError(t, e)
	assert.True(t, b)
	assert.NoError(t, os.RemoveAll(projectDir))
	assert.NoError(t, os.RemoveAll(outputDir))

	// test for -S .tmpl
	inTmplFile, inRegFile, outTmplFile, outRegFile, err = testhelpers.CreateTestTemplateProject(projectDir, subDir, outputDir)
	assert.NoError(t, newProject(projectDir, outputDir, []string{}, []string{testhelpers.TmplSuffix}, []string{}, []string{}))
	// Check that the file suffix was stripped
	assert.NoFileExists(t, outTmplFile)
	assert.FileExists(t, strings.TrimSuffix(outTmplFile, testhelpers.TmplSuffix))
	// Check that outTmplFile didn't get parsed as a template
	b, e = filesAreTheSame(inTmplFile, strings.TrimSuffix(outTmplFile, testhelpers.TmplSuffix))
	assert.NoError(t, e)
	assert.True(t, b)
	assert.NoError(t, os.RemoveAll(projectDir))
	assert.NoError(t, os.RemoveAll(outputDir))

	// test for -t .txt
	inTmplFile, inRegFile, outTmplFile, outRegFile, err = testhelpers.CreateTestTemplateProject(projectDir, subDir, outputDir)
	assert.NoError(t, newProject(projectDir, outputDir, []string{}, []string{}, []string{testhelpers.RegSuffix}, []string{}))
	assert.NoFileExists(t, outRegFile)
	assert.FileExists(t, strings.TrimSuffix(outRegFile, testhelpers.RegSuffix))
	assert.NoError(t, os.RemoveAll(projectDir))
	assert.NoError(t, os.RemoveAll(outputDir))

	// test for -i file.txt
	inTmplFile, inRegFile, outTmplFile, outRegFile, err = testhelpers.CreateTestTemplateProject(projectDir, subDir, outputDir)
	assert.NoError(t, newProject(projectDir, outputDir, []string{}, []string{}, []string{}, []string{path.Base(inRegFile)}))
	assert.NoFileExists(t, outRegFile)
	assert.NoError(t, os.RemoveAll(projectDir))
	assert.NoError(t, os.RemoveAll(outputDir))
}

func testNewProjectNotFound(t *testing.T) {
	assert.Error(t, newProject("", "", []string{}, []string{}, []string{}, []string{}))
}

func testProcessPath(t *testing.T) {
	assert.Error(t, processPath("path", "output", "no", []string{}, []string{}, []string{}, []string{}))
}

func testCopyFile(t *testing.T) {
	testFileA := "testFileA"
	testFileB := "testFileB"

	ioutil.WriteFile(testFileA, []byte("Hello, World"), 0644)
	assert.NoError(t, copyFile(testFileA, testFileB))

	hashOne, err := md5OfFile(testFileA)
	assert.NoError(t, err)
	hashTwo, err := md5OfFile(testFileB)
	assert.NoError(t, err)
	assert.EqualValues(t, hashOne, hashTwo)

	assert.NoError(t, os.Remove(testFileA))
	assert.NoError(t, os.Remove(testFileB))

	assert.Error(t, copyFile("nope", "alsonope"))
}

func testShouldCopyFile(t *testing.T) {
	suffix := ".txt"
	assert.True(t, shouldCopy("file"+suffix, []string{suffix}))
	assert.False(t, shouldCopy("file"+suffix, []string{".noteverathing"}))
}

func testShouldTrimFile(t *testing.T) {
	suffix := ".txt"
	b, s := shouldTrim("file"+suffix, []string{suffix})
	assert.Equal(t, b, true)
	assert.Equal(t, s, suffix)

	b, s = shouldTrim("file"+suffix, []string{".noteverathing"})
	assert.Equal(t, b, false)
	assert.Equal(t, s, "")
}

func testShouldIgnore(t *testing.T) {
	assert.True(t, shouldIgnore("README.md", []string{"1234", "README.md"}))
	assert.False(t, shouldIgnore("README.md", []string{"1234", "abcd"}))
}

func md5OfFile(filePath string) (string, error) {
	var file *os.File
	var err error

	if file, err = os.Open(filePath); err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)[:16]), nil
}

func filesAreTheSame(fileA, fileB string) (bool, error) {
	var hashOne, hashTwo string
	var err error
	if hashOne, err = md5OfFile(fileA); err != nil {
		return false, err
	}
	if hashTwo, err = md5OfFile(fileB); err != nil {
		return false, err
	}
	return hashOne == hashTwo, err
}
