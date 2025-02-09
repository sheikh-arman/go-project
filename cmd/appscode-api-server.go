/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/sheikh-arman/go-project/pkg/appscodeapi/handler"
	"github.com/spf13/cobra"
)

// appscodeApiServerCmd represents the appscodeApiServer command
var appscodeApiServerCmd = &cobra.Command{
	Use:   "appscode-api-server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("appscode-api-erver starting...")
		handler.Handle()
	},
}

func init() {
	rootCmd.AddCommand(appscodeApiServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appscodeApiServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appscodeApiServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
