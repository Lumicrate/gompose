package auth

type AuthUser interface {
	GetID() string
	GetEmail() string
	GetHashedPassword() string
}
