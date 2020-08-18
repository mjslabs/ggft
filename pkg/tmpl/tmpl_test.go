package tmpl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mjslabs/ggft/internal/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestPkgTmpl(t *testing.T) {
	t.Run("CreateFileFromTemplate", testCreateFile)
	t.Run("ScanTemplateForVars", testScanTemplateForVars)
	t.Run("ScanTemplateForVarsInvalid", testScanTemplateForVarsInvalid)
}

func testCreateFile(t *testing.T) {
	projectDir := "testingInPath"
	subDir := "testingSubDir"
	outputDir := "testingOutPath"

	inTmplFile, _, outTmplFile, _, err := testhelpers.CreateTestTemplateProject(projectDir, subDir, outputDir)
	assert.NoError(t, err)
	os.MkdirAll(filepath.Join(outputDir, subDir), 0755)
	assert.NoError(t, CreateFileFromTemplate(inTmplFile, outTmplFile))
	os.RemoveAll(projectDir)
	os.RemoveAll(outputDir)
}

func testScanTemplateForVars(t *testing.T) {
	tmplFile := ".testing.tmpl"
	varName := "TestingVarName"
	assert.NoError(t, testhelpers.CreateFileWithContents(tmplFile, "{{ ."+varName+" }} - {{ .TemplateFileName }}"))
	varList, err := ScanTemplateForVars(tmplFile)

	assert.NoError(t, err)
	assert.Len(t, varList, 1)
	assert.Equal(t, varList[0], varName)
	assert.NoError(t, os.Remove(tmplFile))
}

func testScanTemplateForVarsInvalid(t *testing.T) {
	tmplFile := ".testing.tmpl"
	assert.NoError(t, testhelpers.CreateFileWithContents(tmplFile, "{{ doesntexist.ForSure }} - {{ .TemplateFileName }}"))
	_, err := ScanTemplateForVars(tmplFile)

	assert.Error(t, err)
	assert.NoError(t, os.Remove(tmplFile))
}
