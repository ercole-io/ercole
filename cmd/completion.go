/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion SHELL",
	Short: "Output shell completion code for the specified shell (bash or zsh).",
	Long: `Output shell completion code for the specified shell (bash or zsh).
For zsh, you have to put the output in an _ercole file inside your $fpath.
For bash, you have to put the output in a file parsed by bash, like ~/.bash_completion  
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Shell not specified.")
			return
		}
		if len(args) > 1 {
			fmt.Println("Too many arguments. Expected only the shell type.")
			return
		}

		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
			break
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
			break
		default:
			fmt.Println("Unsupported shell type.")
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
