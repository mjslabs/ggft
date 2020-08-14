/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/mjslabs/ggft/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List downloaded templates available to use with 'new'",
	Run:   cmdList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// List repos that have been downloaded
func cmdList(cmd *cobra.Command, args []string) {
	utils.CheckError(printDirs(viper.GetString("cache")))
}

func printDirs(path string) error {
	dirs, err := listDirs(path)
	if err != nil {
		return err
	}
	fmt.Println(strings.Join(dirs, "\n"))
	return nil
}

func listDirs(path string) (dirList []string, err error) {
	target, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		fmt.Println(target, "doesn't exists")
		return nil, err
	}

	dirs, err := ioutil.ReadDir(target)
	for _, f := range dirs {
		if stat, _ := os.Stat(filepath.Join(target, f.Name())); stat.IsDir() {
			dirList = append(dirList, f.Name())
		}
	}

	return
}
