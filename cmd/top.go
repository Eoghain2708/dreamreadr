/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
	"github.com/spf13/cobra"
)

var validArgs = []string{"emotion", "theme", "symbol", "location", "person"}

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("invalid args")
		}

		var limit int
		if len(args) == 2 {
			l, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			limit = l
		} else {
			limit = 20
		}

		if !slices.Contains(validArgs, args[0]) {
			return fmt.Errorf("invalid arg")
		}

		var result []dream.FeatureCount
		var err error

		switch args[0] {
		case "emotion":
			result, err = dreamService.GetTopEmotions(limit)
		case "symbol":
			result, err = dreamService.GetTopSymbols(limit)
		case "theme":
			result, err = dreamService.GetTopThemes(limit)
		case "person":
			result, err = dreamService.GetTopPeople(limit)
		case "location":
			result, err = dreamService.GetTopLocations(limit)
		}

		if err != nil {
			return err
		}

		for _, fc := range result {
			fmt.Println(fc.Value)
			fmt.Println(fc.Count)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(topCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// topCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// topCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
