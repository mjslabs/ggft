/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/mjslabs/ggft/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <git url> <template name>",
	Short: "Download a template repo and name it for later use",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires a git url and name of template")
		}
		return nil
	},
	Run: cmdGet,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

// Get a template
func cmdGet(cmd *cobra.Command, args []string) {
	utils.CheckError(getTemplate(args[0], filepath.Join(viper.GetString("cache"), args[1])))
}

func getTemplate(gitURL, outputDir string) error {
	target, err := homedir.Expand(outputDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(target); !os.IsNotExist(err) {
		return errors.New(target + " already exists")
	}

	// TODO add logic for auto-naming, e.g. repo ggft-microservice would be named "microservice"
	fmt.Println("Cloning", gitURL, "to", target)
	return utils.GitClone(gitURL, target)
}
