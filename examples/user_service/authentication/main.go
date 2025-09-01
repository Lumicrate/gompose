package main

import (
	"github.com/Lumicrate/gompose/auth/jwt"
	"github.com/Lumicrate/gompose/core"
	"github.com/Lumicrate/gompose/crud"
	"github.com/Lumicrate/gompose/db/postgres"
	"github.com/Lumicrate/gompose/http/gin"
	"strconv"
	"time"
)

// User entity definition with basic fields
type User struct {
	ID          int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

func (u *User) GetID() string {
	return strconv.Itoa(u.ID)
}
func (u *User) GetEmail() string {
	return u.Email
}
func (u *User) GetHashedPassword() string {
	return u.Password
}

type Office struct {
	ID       int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func main() {
	// Configure Postgres DSN
	dsn := "host=localhost user=postgres password=password dbname=mydb port=5432 sslmode=disable"

	// Initialize Postgres DB adapter
	dbAdapter := postgres.New(dsn)

	// Initialize Gin HTTP adapter on port 8080
	httpEngine := ginadapter.New(8080)

	authProvider := jwt.NewJWTAuthProvider("SecretKEY", dbAdapter).SetUserModel(&User{}).SetTokenTTL(time.Hour * 1)

	// Create app instance
	app := core.NewApp().
		AddEntity(Office{}, crud.Protect("POST", "PUT", "DELETE")).
		UseDB(dbAdapter).
		UseHTTP(httpEngine).
		UseAuth(authProvider).
		UseSwagger()
	app.Run()
}
