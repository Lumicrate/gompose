package main

import (
	"github.com/Lumicrate/gompose/auth/jwt"
	"github.com/Lumicrate/gompose/core"
	"github.com/Lumicrate/gompose/crud"
	"github.com/Lumicrate/gompose/db/postgres"
	entities "github.com/Lumicrate/gompose/examples/cmd"
	"github.com/Lumicrate/gompose/http/gin"
)

// sample entity
type User struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	dsn := "host=localhost user=my_user password=my_password dbname=mydb port=5432 sslmode=disable"
	dbAdapter := postgres.New(dsn)
	httpEngine := ginadapter.New(8080)
	authProvider := jwt.NewJWTAuthProvider("SecretKey", dbAdapter)

	app := core.NewApp().
		AddEntity(User{}, crud.Protect("POST", "PUT", "DELETE")).
		AddEntity(entities.Rocket{}, crud.Protect("POST", "PUT", "DELETE")).
		UseDB(dbAdapter).
		UseHTTP(httpEngine).
		UseAuth(authProvider).
		UseSwagger()

	app.Run()
}
