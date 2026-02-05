package cmd

import (
	"fmt"
	"puzzle/generator"

	"github.com/spf13/cobra"
)

var fields string

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Generate auth module",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Fields:", fields)

		generator.GenerateAuth(fields)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.Flags().StringVarP(
		&fields,
		"fields",
		"f",
		"",
		"Auth fields (comma separated)",
	)
}
