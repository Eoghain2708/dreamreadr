/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// occurrencesCmd represents the occurrences command
var occurrencesCmd = &cobra.Command{
	Use:   "occurrences",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var limit int
		var err error
		if len(args) == 0 {
			limit = 100
		} else {
			limit, err = strconv.Atoi(args[0])
		}

		if err != nil {
			return err
		}

		results, err := dreamService.FindCommonCoOccs(limit)
		if err != nil {
			return err
		}

		for _, r := range results {
			fmt.Printf("%s | %s      | %d\n", r.FeatureA, r.FeatureB, r.Count)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(occurrencesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// occurrencesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// occurrencesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
