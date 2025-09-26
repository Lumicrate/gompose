package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Driver string `yaml:"driver"`
		DSN    string `yaml:"dsn"`
		Name   string `yaml:"name"`
	} `yaml:"database"`
	HTTP struct {
		Engine string `yaml:"engine"`
		Port   int    `yaml:"port"`
	} `yaml:"http"`
	Auth struct {
		Secret string `yaml:"secret"`
	} `yaml:"auth"`
}

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new gompose project based on gompose.yaml",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		data, err := os.ReadFile("gompose.yaml")
		if err != nil {
			fmt.Println("gompose.yaml not found. Run 'gompose config' first.")
			return
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			fmt.Println("Failed to parse gompose.yaml:", err)
			return
		}

		if err := os.Mkdir(projectName, 0755); err != nil {
			fmt.Println("Error creating project:", err)
			return
		}

		// Pick template based on db + http
		var mainTemplate string
		if cfg.Database.Driver == "postgres" && cfg.HTTP.Engine == "gin" {
			mainTemplate = postgresGinTemplate
		} else if cfg.Database.Driver == "mongodb" && cfg.HTTP.Engine == "gin" {
			mainTemplate = mongoGinTemplate
		} else {
			fmt.Println("Unsupported db/http combination")
			return
		}

		tmpl, _ := template.New("main").Parse(mainTemplate)
		f, _ := os.Create(fmt.Sprintf("%s/main.go", projectName))
		defer f.Close()
		tmpl.Execute(f, cfg)

		fmt.Printf("Project %s initialized with gompose.yaml configs!\n", projectName)
	},
}

// Templates
const postgresGinTemplate = `package main

import (
    "github.com/Lumicrate/gompose/core"
    "github.com/Lumicrate/gompose/db/postgres"
    "github.com/Lumicrate/gompose/http/gin"
    "github.com/Lumicrate/gompose/auth/jwt"
    "github.com/Lumicrate/gompose/crud"
)

// sample entity
type User struct {
    ID    int    ` + "`json:\"id\" gorm:\"primaryKey;autoIncrement\"`" + `
    Name  string ` + "`json:\"name\"`" + `
    Email string ` + "`json:\"email\"`" + `
}

func main() {
    dsn := "{{.Database.DSN}}"
    dbAdapter := postgres.New(dsn)
    httpEngine := ginadapter.New({{.HTTP.Port}})
    authProvider := jwt.NewJWTAuthProvider("{{.Auth.Secret}}", dbAdapter)

    app := core.NewApp().
        AddEntity(User{}, crud.Protect("POST", "PUT", "DELETE")).
        UseDB(dbAdapter).
        UseHTTP(httpEngine).
        UseAuth(authProvider)

    app.Run()
}
`

const mongoGinTemplate = `package main

import (
    "github.com/Lumicrate/gompose/core"
    "github.com/Lumicrate/gompose/db/mongodb"
	"github.com/Lumicrate/gompose/auth/jwt"
    "github.com/Lumicrate/gompose/http/gin"
	"github.com/Lumicrate/gompose/crud"
)

type User struct {
    ID    int    ` + "`json:\"id\" bson:\"id,omitempty\"`" + `
    Name  string ` + "`json:\"name\" bson:\"name\"`" + `
    Email string ` + "`json:\"email\" bson:\"email\"`" + `
}

func main() {
    mongoURI := "{{.Database.DSN}}"
    dbName := "{{.Database.Name}}"
    dbAdapter := mongodb.New(mongoURI, dbName)
    authProvider := jwt.NewJWTAuthProvider("{{.Auth.Secret}}", dbAdapter)

    httpEngine := ginadapter.New({{.HTTP.Port}})

    app := core.NewApp().
        AddEntity(User{}, crud.Protect("POST", "PUT", "DELETE")).
        UseDB(dbAdapter).
        UseHTTP(httpEngine).
		UseAuth(authProvider)

    app.Run()
}
`
