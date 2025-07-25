package crud

import (
	"github.com/Lumicrate/gompose/auth"
	"github.com/Lumicrate/gompose/db"
	"github.com/Lumicrate/gompose/http"
	"github.com/Lumicrate/gompose/utils"
	"reflect"
	"strings"
)

func RegisterCRUDRoutes(
	engine http.HTTPEngine,
	dbAdapter db.DBAdapter,
	entity any,
	config *Config,
	authProvider auth.AuthProvider,
) {
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	entityName := t.Name()
	basePath := "/" + strings.ToLower(utils.Pluralize(entityName))

	register := func(method, path string, handler http.HandlerFunc) {
		wrapped := handler
		if config.ProtectedMethods[method] && authProvider != nil {
			wrapped = authProvider.Middleware()(handler)
		}
		engine.RegisterRoute(method, path, wrapped)
	}

	// GET /entities (list)
	register("GET", basePath, func(ctx http.Context) {
		handleGetAll(ctx, dbAdapter, entity)
	})

	// GET /entities/:id
	register("GET", basePath+"/:id", func(ctx http.Context) {
		handleGetByID(ctx, dbAdapter, entity)
	})

	// POST /entities
	register("POST", basePath, func(ctx http.Context) {
		handleCreate(ctx, dbAdapter, entity)
	})

	// PUT /entities/:id
	register("PUT", basePath+"/:id", func(ctx http.Context) {
		handleUpdate(ctx, dbAdapter, entity)
	})

	// PATCH /entities/:id
	register("PATCH", basePath+"/:id", func(ctx http.Context) {
		handlePatch(ctx, dbAdapter, entity)
	})

	// DELETE /entities/:id
	register("DELETE", basePath+"/:id", func(ctx http.Context) {
		handleDelete(ctx, dbAdapter, entity)
	})
}
