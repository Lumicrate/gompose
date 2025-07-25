# Changelog

All notable changes to this project will be documented in this file.

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
