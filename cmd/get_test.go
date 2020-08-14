package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mjslabs/ggft/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestCmdGet(t *testing.T) {
	t.Run("getTemplate", testGetTemplate)
	t.Run("getTemplateAlreadyExists", testGetTemplateAlreadyExists)
}

func testGetTemplate(t *testing.T) {
	testDir := "testingProject"
	assert.NoError(t, getTemplate(testhelpers.GitURL, testDir))
	_, err := os.Stat(filepath.Join(testDir, ".git"))
	assert.NoError(t, err)
	assert.NoError(t, os.RemoveAll(testDir))
}

func testGetTemplateAlreadyExists(t *testing.T) {
	testDir := "testingProject"
	assert.NoError(t, os.Mkdir(testDir, 0755))
	assert.Error(t, getTemplate(testhelpers.GitURL, testDir))
	assert.NoError(t, os.RemoveAll(testDir))
}
