package cmd

import (
	"GOTS/internal"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// Command line flags
var srcFile string
var srcFolder string
var outputFile string
var configFile string

var root = &cobra.Command{
	Use:   "gots",
	Short: "pipe the files to convert types from",
	Long:  `pipe the files to convert types from`,
	// TODO : make a function to that the file can't be inserted with the folder
	Run: func(cmd *cobra.Command, args []string) {
		var ctx context.Context = context.Background()
		if srcFile == "" && srcFolder == "" {
			fmt.Println("No source file or folder provided, check the help command")
			return
		}
		if srcFile != "" && srcFolder != "" {
			fmt.Println("You can't provide both a file and a folder, just one")
			return
		}

		if configFile != "" {
			config := internal.NewConfig(configFile)
			config.LoadContent()
			config.ParseContent()

			ctx = context.WithValue(ctx, "config", config)
		}

		switch true {
		case srcFile != "":
			internal.Transpile(srcFile, outputFile, ctx.Value("config").(*internal.Config))
		case srcFolder != "":
			fmt.Println("Folder functionality, Coming soon...")
		default:
			fmt.Println("No source file or folder provided")
		}
	},
}

func Execute() error {
	root.Flags().StringVarP(&srcFile, "src", "s", "", "source file to convert types from")
	root.Flags().StringVarP(&srcFolder, "folder", "f", "", "source folder to convert types from")
	root.Flags().StringVarP(&outputFile, "output", "o", "", "output file name")
	root.Flags().StringVarP(&configFile, "config", "c", "", "config file to use")
	root.MarkFlagRequired("src")
	root.MarkFlagRequired("output")

	return root.Execute()
}
