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

	"github.com/go-git/go-git/v5"
	"github.com/mitchellh/go-homedir"
	"github.com/mjslabs/ggft/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Attempt to update all templates that have been downloaded",
	Run:   cmdUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// Update all templates from origin
func cmdUpdate(cmd *cobra.Command, args []string) {
	utils.CheckError(updateProjects(filepath.Join(viper.GetString("cache"))))
}

func updateProjects(projectsPath string) error {
	templateDir, err := homedir.Expand(projectsPath)
	if err != nil {
		return err
	}

	dirContents, err := ioutil.ReadDir(templateDir)
	for _, f := range dirContents {
		fullPath := filepath.Join(templateDir, f.Name())
		if _, err := os.Stat(filepath.Join(fullPath, ".git")); os.IsNotExist(err) {
			continue
		}
		if err := updateRepo(fullPath); err != nil {
			return err
		}
	}

	return nil
}

func updateRepo(path string) error {
	var (
		r   *git.Repository
		w   *git.Worktree
		err error
	)

	fmt.Println("updating", path)

	if r, err = git.PlainOpen(path); err != nil {
		fmt.Println("can't open", path)
		return err
	}

	if w, err = r.Worktree(); err != nil {
		return err
	}
	if err = w.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil && !strings.Contains(err.Error(), "already up-to-date") {
		fmt.Println("error pulling repo", path, ":", err.Error())
		return err
	}

	return nil
}
