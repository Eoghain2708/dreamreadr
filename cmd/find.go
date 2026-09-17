/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
	"github.com/spf13/cobra"
)

var (
	findEmotion  string
	findSymbol   string
	findLocation string
	findPerson   string
	findTheme    string
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var filters []dream.FeatureFilter
		if findEmotion != "" {
			filters = append(filters, dream.FeatureFilter{
				Type:  dream.Emotion,
				Value: findEmotion,
			})
		}
		if findLocation != "" {
			filters = append(filters, dream.FeatureFilter{
				Type:  dream.Location,
				Value: findLocation,
			})
		}
		if findPerson != "" {
			filters = append(filters, dream.FeatureFilter{
				Type:  dream.Person,
				Value: findPerson,
			})
		}
		if findTheme != "" {
			filters = append(filters, dream.FeatureFilter{
				Type:  dream.Theme,
				Value: findTheme,
			})
		}
		if findSymbol != "" {
			filters = append(filters, dream.FeatureFilter{
				Type:  dream.Symbol,
				Value: findSymbol,
			})
		}

		dreams, err := dreamService.FindDreams(filters...)
		if err != nil {
			return err
		}

		for _, dream := range dreams {
			fmt.Println(dream.Title)
			fmt.Println(dream.ID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().StringVar(
		&findEmotion, "emotion", "", "match emotion",
	)
	findCmd.Flags().StringVar(
		&findSymbol, "symbol", "", "match symbol",
	)
	findCmd.Flags().StringVar(
		&findLocation, "location", "", "match location",
	)
	findCmd.Flags().StringVar(
		&findTheme, "theme", "", "match theme",
	)
	findCmd.Flags().StringVar(
		&findPerson, "person", "", "match person",
	)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
