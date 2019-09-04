package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init initialize a git repository",
	Long:  "init initialize a git repository with default content",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing git add ...")
	},
}
