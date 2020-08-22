package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mjslabs/ggft/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestCmdNew(t *testing.T) {
	t.Run("NewProject", testNewProject)
	t.Run("NewProject", testNewProjectNotFound)
	t.Run("CopyFile", testCopyFile)
	t.Run("ShouldCopyFile", testShouldCopyFile)
	t.Run("ShouldCopyAndTrimFile", testShouldCopyAndTrimFile)
}

func testNewProject(t *testing.T) {
	projectDir := "testingInPath"
	subDir := "testingSubDir"
	outputDir := "testingOutPath"

	inTmplFile, inRegFile, outTmplFile, outRegFile, err := testhelpers.CreateTestTemplateProject(projectDir, subDir, outputDir)
	assert.NoError(t, err)
	assert.NoError(t, newProject(projectDir, outputDir, []string{}, []string{}))

	hashOne, err := md5OfFile(inTmplFile)
	assert.NoError(t, err)
	hashTwo, err := md5OfFile(outTmplFile)
	assert.NoError(t, err)
	assert.NotEqual(t, hashOne, hashTwo)

	hashOne, err = md5OfFile(inRegFile)
	assert.NoError(t, err)
	hashTwo, err = md5OfFile(outRegFile)
	assert.NoError(t, err)
	assert.Equal(t, hashOne, hashTwo)

	assert.NoError(t, os.RemoveAll(projectDir))
	assert.NoError(t, os.RemoveAll(outputDir))
}

func testNewProjectNotFound(t *testing.T) {
	assert.Error(t, newProject("", "", []string{}, []string{}))
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

func testShouldCopyAndTrimFile(t *testing.T) {
	suffix := ".txt"
	b, s := shouldCopyAndTrim("file"+suffix, []string{suffix})
	assert.Equal(t, b, true)
	assert.Equal(t, s, suffix)

	b, s = shouldCopyAndTrim("file"+suffix, []string{".noteverathing"})
	assert.Equal(t, b, false)
	assert.Equal(t, s, "")
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
