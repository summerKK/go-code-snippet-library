package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  "",
}

func init() {
	rootCmd.AddCommand(wordCmd)
	rootCmd.AddCommand(timeCmd)
	rootCmd.AddCommand(sqlCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
