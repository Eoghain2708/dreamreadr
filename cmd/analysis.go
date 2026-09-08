/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// analysisCmd represents the analysis command
var analysisCmd = &cobra.Command{
	Use:   "analysis",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		analysis, err := dreamService.GetDreamAnalysis(args[0])
		if err != nil {
			return nil
		}

		fmt.Printf("%v", analysis)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(analysisCmd)

}
