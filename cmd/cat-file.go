package cmd

import (
	"fmt"
	"github.com/ajanthan/gitcmd/pkg"
	"github.com/ajanthan/gitcmd/pkg/objects"
	"github.com/spf13/cobra"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(catFileCommand)
}

var catFileCommand = &cobra.Command{
	Use:   "cat-file [type] [object]",
	Short: "out put content of object",
	Args:  cobra.MinimumNArgs(1),
	Long:  "Find a object by type and hash and display in the console",
	Run: func(cmd *cobra.Command, args []string) {
		var objectType, object string
		if len(args) == 1 {
			objectType = "blob"
			object = args[0]
		} else {
			objectType = args[0]
			object = args[1]
		}
		workingDirectory, err := pkg.FindRepository(filepath.Clean("."))
		if err != nil {
			fmt.Println(err)
			return
		}
		gitObject, err := objects.DecodeGitObject(object, objectType, workingDirectory)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s", gitObject.Data)
	},
}
