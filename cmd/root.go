/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"fmt"
	"os"

	"github.com/mjslabs/ggft/pkg/config"
	"github.com/spf13/cobra"
)

var version = "undefined"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ggft",
	Short: "Generate a new Go project from a given framework template",
	Long: `Download, list, and delete templates. Create new directories based on a starting template. For example:

ggft get https://github.com/mjslabs/ggft-template service
ggft list
ggft new service my-new-service
`,

	Version: version,
}

func init() {
	cobra.OnInitialize(func() {
		config.InitAndLoadConfig()
	})
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
