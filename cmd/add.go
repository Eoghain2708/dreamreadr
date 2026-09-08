/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Title: > ")
		title, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		title = strings.TrimSpace(title)

		var content strings.Builder
		fmt.Println("Content: ")
		fmt.Println("Type END on its own line when finished")
		fmt.Print("> ")
		for {
			line, _ := reader.ReadString('\n')

			if strings.TrimSpace(line) == "END" {
				break
			}
			content.WriteString(line)
		}

		dream, err := dreamService.CreateDream(title, content.String())
		if err != nil {
			return fmt.Errorf("error creating dream: %v", err)
		}
		fmt.Printf("Dream created! ID: %s\n", dream.ID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
