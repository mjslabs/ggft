package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mjslabs/ggft/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestCmdUpdate(t *testing.T) {
	t.Run("updateProjects", testUpdateProjects)
}

func testUpdateProjects(t *testing.T) {
	testDir := "testingParentDir"
	testProjDir := filepath.Join(testDir, "testingProject")

	os.Mkdir(testDir, 0755)
	assert.NoError(t, getTemplate(testhelpers.GitURL, testProjDir))
	_, err := os.Stat(filepath.Join(testProjDir, ".git"))
	assert.NoError(t, err)
	assert.NoError(t, updateProjects(testDir))
	assert.NoError(t, os.RemoveAll(testDir))

	assert.Error(t, updateRepo("nope"))
}
