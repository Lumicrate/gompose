package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dbFlag, httpFlag, dsnFlag, secretFlag, dbNameFlag string
var portFlag int

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate a gompose.yaml config file",
	Run: func(cmd *cobra.Command, args []string) {
		config := fmt.Sprintf(`database:
  driver: %s
  dsn: "%s"
  name: %s

http:
  engine: %s
  port: %d

auth:
  secret: "%s"
`, dbFlag, dsnFlag, dbNameFlag, httpFlag, portFlag, secretFlag)

		if err := os.WriteFile("gompose.yaml", []byte(config), 0644); err != nil {
			fmt.Println("Error creating gompose.yaml:", err)
			return
		}

		fmt.Println("gompose.yaml generated successfully with custom settings.")
	},
}

func init() {
	configCmd.Flags().StringVar(&dbFlag, "db", "postgres", "Database driver (postgres|mongodb)")
	configCmd.Flags().StringVar(&httpFlag, "http", "gin", "HTTP engine (gin)")
	configCmd.Flags().StringVar(&dsnFlag, "dsn", "host=localhost user=username password=password dbname=mydb port=5432 sslmode=disable", "Database DSN/URI")
	configCmd.Flags().StringVar(&dbNameFlag, "dbname", "mydb", "Database name (used by MongoDB)")
	configCmd.Flags().IntVar(&portFlag, "port", 8080, "HTTP port")
	configCmd.Flags().StringVar(&secretFlag, "secret", "SecretKEY", "Auth secret")
}
