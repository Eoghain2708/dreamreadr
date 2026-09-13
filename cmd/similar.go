/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// similarCmd represents the similar command
var similarCmd = &cobra.Command{
	Use:   "similar",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		ctx := context.Background()
		dreams, err := dreamService.FindSimilarDreams(ctx, args[0], limit)
		if err != nil {
			return err
		}

		for _, d := range dreams {
			fmt.Printf("Similarity: %.2f - title: %s\n", d.Similarity, d.Dream.Title)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(similarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// similarCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// similarCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
