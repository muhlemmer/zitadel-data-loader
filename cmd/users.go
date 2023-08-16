/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("users called")
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
}
