/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"errors"
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
	newCmd.Flags().StringSliceP("ignore", "i", []string{"README.md"}, "Files to ignore")
	newCmd.Flags().StringSliceP("skip", "s", []string{}, "Suffixes of files to copy instead of parse")
	newCmd.Flags().StringSliceP("skiptrim", "S", []string{}, "Suffixes of files to copy instead of parse, trimming suffix")
	newCmd.Flags().StringSliceP("trim", "t", []string{".ggft"}, "Suffixes of files to parse, trimming suffix")
}

// Create a new directory from a template
func cmdNew(cmd *cobra.Command, args []string) {
	skipSuffixes, err := cmd.Flags().GetStringSlice("skip")
	utils.CheckError(err)
	skipTrimSuffixes, err := cmd.Flags().GetStringSlice("skiptrim")
	utils.CheckError(err)
	trimSuffixes, err := cmd.Flags().GetStringSlice("trim")
	utils.CheckError(err)
	ignoredFiles, err := cmd.Flags().GetStringSlice("ignore")
	utils.CheckError(err)
	utils.CheckError(newProject(filepath.Join(viper.GetString("cache"), args[0]), args[1], skipSuffixes, skipTrimSuffixes, trimSuffixes, ignoredFiles))
}

func newProject(templatePath, outputDir string, skipSuffixes, skipTrimSuffixes, trimSuffixes, ignoredFiles []string) error {
	templateDir, err := homedir.Expand(templatePath)
	if err != nil {
		return err
	}

	// Check for template dir existance
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return errors.New("directory not found: " + templateDir)
	}

	// Create output dir
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}

	// Get template version
	tmpl.SetVar("GitHash", utils.GetGitHash(templateDir))

	// Convert template to output
	return walkTemplateDirectory(templateDir, outputDir, skipSuffixes, skipTrimSuffixes, trimSuffixes, ignoredFiles)
}

func walkTemplateDirectory(templateDir, outputDir string, skipSuffixes, skipTrimSuffixes, trimSuffixes, ignoredFiles []string) error {
	if err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		return processPath(path, templateDir, outputDir, skipSuffixes, skipTrimSuffixes, trimSuffixes, ignoredFiles)
	}); err != nil {
		return err
	}
	return nil
}

func processPath(curPath, templateDir, output string, skipSuffixes, skipTrimSuffixes, trimSuffixes, ignoredFiles []string) error {
	destination := filepath.Join(output, strings.Replace(curPath, templateDir, "", 1))
	fi, err := os.Stat(curPath)
	if err != nil {
		return err
	}

	switch mode := fi.Mode(); {
	case strings.Contains(curPath, string(filepath.Separator)+".git"+string(filepath.Separator)):
		return nil

	case shouldIgnore(curPath, ignoredFiles):
		// Should be ignored because of -i
		return nil

	case mode.IsDir():
		os.Mkdir(destination, 0755)
		return nil

	case shouldCopy(curPath, append(skipSuffixes, skipTrimSuffixes...)):
		// Either -s or -S brought us here
		if b, s := shouldTrim(curPath, skipTrimSuffixes); b {
			return copyFile(curPath, strings.TrimRight(destination, s))
		}
		return copyFile(curPath, destination)

	case mode.IsRegular():
		// Template file, process accordingly
		if err := getVarsFromUser(curPath); err != nil {
			return err
		}
		// Check for -t
		if b, s := shouldTrim(curPath, trimSuffixes); b {
			destination = strings.TrimRight(destination, s)
		}
		return tmpl.CreateFileFromTemplate(curPath, destination)
	}

	return nil
}

func shouldIgnore(path string, ignoredFiles []string) bool {
	for _, s := range ignoredFiles {
		if strings.HasSuffix(path, s) {
			return true
		}
	}
	return false
}

func shouldTrim(path string, trimSuffixes []string) (bool, string) {
	for _, s := range trimSuffixes {
		if strings.HasSuffix(path, s) {
			return true, s
		}
	}
	return false, ""
}

func shouldCopy(path string, skipSuffixes []string) bool {
	for _, s := range skipSuffixes {
		if strings.HasSuffix(path, s) {
			return true
		}
	}
	return false
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
