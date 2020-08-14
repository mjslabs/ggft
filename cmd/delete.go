/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/mjslabs/ggft/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <template>",
	Short: "Delete a template",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a template name")
		}
		return nil
	},
	Run: cmdDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

// Delete a template
func cmdDelete(cmd *cobra.Command, args []string) {
	utils.CheckError(deleteTemplate(filepath.Join(viper.GetString("cache"), args[0])))
}

func deleteTemplate(dir string) error {
	target, err := homedir.Expand(dir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		return err
	}

	return os.RemoveAll(target)
}
