package cmd

import (
	"fmt"
	"gots/internal"

	"github.com/spf13/cobra"
)

// Command line flags
var srcFile string
var srcFolder string
var outputFile string
var configFile string
var recursiveTranspile bool

var root = &cobra.Command{
	Use:   "gots",
	Short: "pipe the files to convert types from",
	Long:  `pipe the files to convert types from`,
	Run: func(cmd *cobra.Command, args []string) {
		// * load the config file
		config := internal.NewConfig(configFile)
		config.ParseContent()

		if srcFile == "" && srcFolder == "" {
			fmt.Println("No source file or folder provided, check the help command")
			return
		}
		if srcFile != "" && srcFolder != "" {
			fmt.Println("You can't provide both a file and a folder, just one")
			return
		}

		switch true {
		case srcFile != "":
			transpiler := internal.NewTranspiler(srcFile, outputFile, config)
			content := transpiler.Transpile()
			internal.CreateFile(&outputFile, content)
		case srcFolder != "":
			Files := internal.MultiFile(srcFolder, recursiveTranspile)
			for i, file := range Files {
				transpiler := internal.NewTranspiler(file, outputFile, config)
				content := transpiler.Transpile()
				if i == 0 {
					internal.CreateFile(&outputFile, content)
				} else {
					internal.UpdateFile(&outputFile, content)
				}
			}
		default:
			fmt.Println("No source file or folder provided")
		}
	},
}

func Execute() error {
	root.Flags().StringVarP(&srcFile, "src", "s", "", "source file to convert types from")
	root.Flags().StringVarP(&srcFolder, "folder", "f", "", "source folder to convert types from")
	root.Flags().StringVarP(&outputFile, "output", "o", "", "output file name")
	root.Flags().StringVarP(&configFile, "config", "c", "./gots.yaml", "config file to use")
	root.Flags().BoolVarP(&recursiveTranspile, "recursive", "r", false, "recursively transpile all nested files in the folder")
	root.MarkFlagRequired("output")

	return root.Execute()
}
