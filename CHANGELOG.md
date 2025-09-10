# Changelog

All notable changes to this project will be documented in this file.

## [v1.3.0] - 2025-09-10
### Added
- **i18n / Translator support**:
  - Ability to internationalize and translate texts using JSON/YAML files in the `locales` directory.
  - Support for multiple languages (e.g., `fa`, `en`) with a configurable default language.
  - Support for dynamic parameters in translation strings.

## [v1.2.1] - 2025-09-01
### Added
- Support for **query parameters** in GET routes:
  - `limit` and `offset` for pagination
  - `sort` for sorting (e.g., `sort=name,-created_at`)
  - Entity fields as **filters** (e.g., `?name=john`)

### Changed
- Refactored `SwaggerProvider.Generate` for clarity and extensibility
- Added detailed inline comments in Swagger generator code

### Fixed
- Missing documentation for query-based filters in GET endpoints

## [1.2.0] - 2025-08-24
### Added
- **Swagger Integration**:
  - Introduced `SwaggerProvider` to auto-generate OpenAPI 3.0 documentation.
  - `/swagger.json` endpoint for raw JSON spec.
  - `/swagger-ui` endpoint with modern Swagger UI, minimal and formal style.
  - Auto-generates request/response schemas from Go structs.
  - Automatically maps `:id` path parameters to `{id}` in Swagger paths.
  - Displays JWT-protected endpoints and allows authentication via Swagger UI.
- **Gin Adapter** updates:
  - `HttpEngine Context` now supports `Body()` to serve raw HTML for Swagger UI.
  - `RegisterRoute()` accepts an optional entity for schema generation.
- **JWT / Auth Enhancements**:
  - Swagger now recognizes endpoints protected by JWT (via `crud.Protect()` or `ProtectAll()`).
  - Swagger UI allows “Authorize” to input JWT tokens for secured endpoints.

## [1.1.0] - 2025-07-25
### Added
- Introduced the `/auth` package with an `AuthUser` interface and a default `UserModel`.
- Added `JWTAuthProvider` under `/auth/jwt`:
  - Auto-registers `/auth/login` and `/auth/register` routes.
  - Includes built-in middleware for JWT-based route protection.
  - Supports `SetTokenTTL()` for customizable token expiration.
  - Allows injecting a custom user model via `SetUserModel()`.
- Added `/crud` helpers:
  - `Protect()` to secure specific HTTP methods.
  - `ProtectAll()` to secure all methods on an entity.
- Added `/utils` helpers:
  - JWT generation/validation.
  - Password hashing/comparison using `bcrypt`.
  - UUID generation.
  - Bearer token extraction from headers.
### Fixed
- Updated MongoDB collection naming to use `utils.Pluralize` for proper pluralization instead of simply appending `"s"`.
- Updated middleware interface to use `http.MiddlewareFunc` with standard `next http.HandlerFunc` chaining for proper execution flow.

## [1.0.1] - 2025-07-18
### Fixed
- Corrected the `basePath` generation to use proper pluralization via the `utils.Pluralize` function instead of manually appending `"s"`.  
  This ensures correct plural forms for entity names (e.g., "category" → "categories").
