package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdDelete(t *testing.T) {
	t.Run("DeleteTemplate", testDeleteTemplate)
}

func testDeleteTemplate(t *testing.T) {
	testDir := "testingDir"
	testSubDirs := []string{"dir1", "dir2", "dir3"}
	for d := range testSubDirs {
		assert.NoError(t, os.MkdirAll(filepath.Join(testDir, testSubDirs[d]), 0755))
	}
	assert.NoError(t, deleteTemplate(testDir))
	_, err := os.Stat(testDir)
	assert.True(t, os.IsNotExist(err))

	assert.Error(t, deleteTemplate("nope"))
}
