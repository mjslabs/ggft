package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdList(t *testing.T) {
	t.Run("ListDirs", testListDirs)
}

func testListDirs(t *testing.T) {
	testDir := "testingDir"
	testSubDirs := []string{"dir1", "dir2", "dir3"}
	for d := range testSubDirs {
		assert.NoError(t, os.MkdirAll(filepath.Join(testDir, testSubDirs[d]), 0755))
	}
	dirs, err := listDirs(testDir)
	assert.NoError(t, err)
	assert.Equal(t, len(testSubDirs), len(dirs))
	assert.NoError(t, os.RemoveAll(testDir))

	_, err = listDirs("nope")
	assert.Error(t, err)
}
