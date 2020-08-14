/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mjslabs/ggft/pkg/askuser"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/mjslabs/ggft/pkg/tmpl"
	"github.com/mjslabs/ggft/pkg/utils"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <type> <name of directory to create>",
	Short: "Generate directory from template",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires a project type and output directory name to create")
		}
		return nil
	},
	Run: cmdNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

// Create a new directory from a template
func cmdNew(cmd *cobra.Command, args []string) {
	utils.CheckError(newProject(filepath.Join(viper.GetString("cache"), args[0]), args[1]))
}

func newProject(templatePath, outputDir string) error {
	templateDir, err := homedir.Expand(templatePath)
	if err != nil {
		return err
	}

	// Check for template dir existance
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		fmt.Println("directory not found:", templateDir)
		os.Exit(1)
	}

	// Create output dir
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}

	// Get template version
	tmpl.SetVar("GitHash", utils.GetGitHash(templateDir))

	// Convert template to output
	return walkTemplateDirectory(templateDir, outputDir)
}

func walkTemplateDirectory(templateDir, outputDir string) error {
	if err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		return processPath(path, templateDir, outputDir)
	}); err != nil {
		return err
	}
	return nil
}

func processPath(curPath, templateDir, output string) error {
	destination := filepath.Join(output, strings.Replace(curPath, templateDir, "", 1))
	fi, err := os.Stat(curPath)
	if err != nil {
		return err
	}

	switch mode := fi.Mode(); {
	case strings.Contains(curPath, string(filepath.Separator)+".git"+string(filepath.Separator)):
		return nil
	case mode.IsDir():
		os.Mkdir(destination, 0755)
		return nil
	case strings.HasSuffix(curPath, ".tmpl"):
		copyFile(curPath, destination)
		return nil
	case mode.IsRegular():
		// Template file, process accordingly
		getVarsFromUser(curPath)
		return tmpl.CreateFileFromTemplate(curPath, destination)
	}

	return nil
}

func getVarsFromUser(tmplPath string) error {
	var varList []string
	var err error

	if varList, err = tmpl.ScanTemplateForVars(string(tmplPath)); err != nil {
		return err
	}

	for i := range varList {
		tmpl.SetVar(varList[i], askuser.Terminal("Enter a value for "+varList[i]+": ", ""))
	}

	return nil
}

// copyFile -
func copyFile(input, output string) error {
	var (
		from, to *os.File
		err      error
	)

	if from, err = os.Open(input); err != nil {
		return err
	}
	defer from.Close()

	if to, err = os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		return err
	}
	defer to.Close()

	if _, err = io.Copy(to, from); err != nil {
		return err
	}

	return nil
}
