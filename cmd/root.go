package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "achim",
	Short: "Advanced Cloud Hyperscaling Infrastructure Manager",
	Long:  "achim manages Exoscale compute instances for education",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Usage: achim [noun] [verb] (or use --help)")
	},
}

var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "manage API access",
	Long:  "manage API access tokens for Exoscale",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: achim access")
	},
}

var accessAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add API access",
	Long:  "add API access for Exoscale",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: achim access add")
	},
}

func init() {
	rootCmd.AddCommand(accessCmd)
	accessCmd.AddCommand(accessAddCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
