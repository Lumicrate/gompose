package core

import (
	"github.com/Lumicrate/gompose/auth"
	"github.com/Lumicrate/gompose/crud"
	"github.com/Lumicrate/gompose/db"
	"github.com/Lumicrate/gompose/http"
	"log"
)

type App struct {
	entities     []registeredEntity
	dbAdapter    db.DBAdapter
	httpEngine   http.HTTPEngine
	middlewares  []http.MiddlewareFunc
	authProvider auth.AuthProvider
}

type registeredEntity struct {
	entity any
	config *crud.Config
}

func NewApp() *App {
	return &App{
		entities:    []registeredEntity{},
		middlewares: []http.MiddlewareFunc{},
	}
}

func (a *App) Entities() []any {
	raw := make([]any, len(a.entities))
	for i, e := range a.entities {
		raw[i] = e.entity
	}
	return raw
}

func (a *App) AddEntity(entity any, opts ...crud.Option) *App {
	cfg := crud.DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	a.entities = append(a.entities, registeredEntity{entity: entity, config: cfg})
	return a
}

func (a *App) UseDB(adapter db.DBAdapter) *App {
	a.dbAdapter = adapter
	return a
}

func (a *App) UseHTTP(engine http.HTTPEngine) *App {
	a.httpEngine = engine
	return a
}

func (a *App) RegisterMiddleware(m http.MiddlewareFunc) *App {
	a.middlewares = append(a.middlewares, m)
	return a
}

func (a *App) UseAuth(provider auth.AuthProvider) *App {
	a.authProvider = provider
	return a
}

func (a *App) Run() {
	if err := a.dbAdapter.Init(); err != nil {
		log.Fatalf("DB Init failed: %v", err)
	}

	if err := a.dbAdapter.Migrate(a.Entities()); err != nil {
		log.Fatalf("DB Migration failed: %v", err)
	}

	if a.authProvider != nil {
		err := a.authProvider.Init()
		if err != nil {
			return
		}
		a.authProvider.RegisterRoutes(a.httpEngine)
	}

	for _, m := range a.middlewares {
		a.httpEngine.Use(m)
	}

	for _, e := range a.entities {
		crud.RegisterCRUDRoutes(a.httpEngine, a.dbAdapter, e.entity, e.config, a.authProvider)
	}

	if err := a.httpEngine.Start(); err != nil {
		log.Fatalf("HTTP Server failed: %v", err)
	}
}
