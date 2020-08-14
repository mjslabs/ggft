package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mjslabs/ggft/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestPkgUtils(t *testing.T) {
	t.Run("GitClone", testGitClone)
	t.Run("GetGitHash", testGitHash)
}

func testGitClone(t *testing.T) {
	testDir := "testingGitClone"
	assert.NoError(t, GitClone(testhelpers.GitURL, testDir))
	assert.DirExists(t, testDir)
	assert.DirExists(t, filepath.Join(testDir, ".git"))
	assert.NoError(t, os.RemoveAll(testDir))
}

func testGitHash(t *testing.T) {
	testDir := "testingGitClone"
	assert.NoError(t, GitClone(testhelpers.GitURL, testDir))
	assert.NotEmpty(t, GetGitHash(testDir))
	assert.NoError(t, os.RemoveAll(testDir))
}
