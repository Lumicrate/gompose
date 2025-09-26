package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate entity [name] --fields name:type,email:type",
	Short: "Generate a new entity struct",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entityName := args[0]
		fields, _ := cmd.Flags().GetString("fields")
		fieldList := strings.Split(fields, ",")

		var fieldStr string
		for _, f := range fieldList {
			parts := strings.Split(f, ":")
			if len(parts) == 2 {
				fieldStr += fmt.Sprintf("    %s %s `json:\"%s\"`\n", strings.Title(parts[0]), strings.Title(parts[1]), parts[0])
			}
		}

		tmplText := `package entities

type {{.Name}} struct {
    ID int ` + "`json:\"id\" gorm:\"primaryKey;autoIncrement\"`" + `
{{.Fields}}}
`
		tmpl, _ := template.New("entity").Parse(tmplText)
		f, _ := os.Create(fmt.Sprintf("%s.go", strings.ToLower(entityName)))
		defer f.Close()
		tmpl.Execute(f, map[string]string{"Name": entityName, "Fields": fieldStr})

		fmt.Printf("Entity %s generated!\n", entityName)
	},
}

func init() {
	generateCmd.Flags().String("fields", "", "Comma separated list of fields (e.g. name:string,email:string)")
}
