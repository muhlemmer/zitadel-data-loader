/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	usersv1 "github.com/muhlemmer/zitadel-data-loader/internal/client/users/users/v1"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create random users",
	Long: `Create users with random fields, such as name, email and profile details.
	The password is set as a hash of 'Password1!'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return usersv1.ImportUsers(background, clientConn, int(userCreateNumber))
	},
}

var (
	userCreateNumber uint
)

func init() {
	usersCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().UintVarP(&userCreateNumber, "number", "n", 100, "Number of users to create")
}
