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
	t.Run("CopyFile", testCopyFile)
}

func testNewProject(t *testing.T) {
	projectDir := "testingInPath"
	subDir := "testingSubDir"
	outputDir := "testingOutPath"

	inTmplFile, inRegFile, outTmplFile, outRegFile, err := testhelpers.CreateTestTemplateProject(projectDir, subDir, outputDir)
	assert.NoError(t, err)
	assert.NoError(t, newProject(projectDir, outputDir))

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

	os.RemoveAll(projectDir)
	os.RemoveAll(outputDir)
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

	os.Remove(testFileA)
	os.Remove(testFileB)
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
