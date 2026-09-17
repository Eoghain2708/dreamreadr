/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// overlapCmd represents the overlap command
var overlapCmd = &cobra.Command{
	Use:   "overlap",
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

		str := strings.Join(args, " ")
		q := strings.Split(str, ", ")

		dreams, err := dreamService.FindOverlappingFeatures(q)
		if err != nil {
			return err
		}

		for _, d := range dreams {
			fmt.Println(d.Title)
			fmt.Println(d.ID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(overlapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// overlapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// overlapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
