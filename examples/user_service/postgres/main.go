package main

import (
	"errors"
	"github.com/Lumicrate/gompose/core"
	"github.com/Lumicrate/gompose/db/postgres"
	"github.com/Lumicrate/gompose/http"
	"github.com/Lumicrate/gompose/http/gin"
	"github.com/Lumicrate/gompose/http/middlewares"
	"log"
	"strings"
	"time"
)

// User entity definition with basic fields
type User struct {
	ID    int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Office struct {
	ID       int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func (u *User) BeforeCreate() error {
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func (u *User) AfterDelete() error {
	log.Printf("User %d deleted", u.ID)
	return nil
}

// Make you own custom middlewares
func CORSMiddleware() http.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(ctx http.Context) {
			ctx.SetHeader("Access-Control-Allow-Origin", "*")
			ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.SetHeader("Access-Control-Allow-Headers", "Authorization, Content-Type")

			if ctx.Method() == "OPTIONS" {
				ctx.SetStatus(204)
				ctx.Abort()
				return
			}
			ctx.Next()
		}
	}
}

func main() {
	// Configure Postgres DSN
	dsn := "host=localhost user=username password=password dbname=mydb port=5432 sslmode=disable"

	// Initialize Postgres DB adapter
	dbAdapter := postgres.New(dsn)

	// Initialize Gin HTTP adapter on port 8080
	httpEngine := ginadapter.New(8080)

	// Create app instance
	app := core.NewApp().
		AddEntity(User{}).
		AddEntity(Office{}).
		UseDB(dbAdapter).
		UseHTTP(httpEngine).
		RegisterMiddleware(middlewares.LoggingMiddleware()).
		RegisterMiddleware(middlewares.RateLimitMiddleware(1 * time.Second)).
		RegisterMiddleware(CORSMiddleware())

	// Run the app
	log.Println("Starting service on :8080")
	app.Run()
}
