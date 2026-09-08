/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// analyseCmd represents the analyse command
var analyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		analysis, err := dreamService.AnalyseDream(ctx, args[0])
		if err != nil {
			return err
		}

		fmt.Printf("%v", analysis)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// analyseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// analyseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
