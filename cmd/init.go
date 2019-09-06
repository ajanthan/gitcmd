package cmd

import (
	"fmt"
	"github.com/ajanthan/gitcmd/pkg"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [repository to create]",
	Short: "init initialize a git repository",
	Args:  cobra.MinimumNArgs(1),
	Long:  "init initialize a git repository with default content",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Executing git add ...%s", args[0])
		repository, err := pkg.NewRepository(args[0], true)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = repository.Create()
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
