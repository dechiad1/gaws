package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

//defined here - used in sg, ec2 subcommands
var (
	sess *session.Session
)

var rootCmd = &cobra.Command{
	Use:   "gaws",
	Short: "gaws is a command line utility to access formatted aws data",
	Long:  `A command line utility to view useful aws data related to ec2 and vpc services`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("usage: gaws <ec2 | sg> [options]")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
