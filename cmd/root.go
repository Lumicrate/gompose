package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gompose",
	Short: "Gompose - Rapid service scaffolding in Go",
	Long:  `Gompose helps you spin up services quickly with DB, HTTP, Auth, and CRUD scaffolding.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'gompose --help' to see available commands")
	},
}

// Execute adds all child commands to the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(generateCmd)
}
