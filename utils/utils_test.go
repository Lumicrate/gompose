package utils_test

import (
	"testing"
	"time"

	"github.com/Lumicrate/gompose/utils"
	"github.com/stretchr/testify/require"
	"regexp"
	"strings"
)

// ExtractBearerToken

func TestExtractBearerToken_Valid(t *testing.T) {
	token, err := utils.ExtractBearerToken("Bearer abc123")
	require.NoError(t, err)
	require.Equal(t, "abc123", token)
}

func TestExtractBearerToken_MissingHeader(t *testing.T) {
	_, err := utils.ExtractBearerToken("")
	require.Error(t, err)
}

func TestExtractBearerToken_InvalidHeader(t *testing.T) {
	_, err := utils.ExtractBearerToken("Token abc123")
	require.Error(t, err)
}

// Password Hashing

func TestGenerateFromPassword_ProducesHash(t *testing.T) {
	hash, err := utils.GenerateFromPassword("mypassword")
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEqual(t, "mypassword", hash)
}

func TestCompareHashAndPassword_CorrectPassword(t *testing.T) {
	hash, err := utils.GenerateFromPassword("secret")
	require.NoError(t, err)

	err = utils.CompareHashAndPassword(hash, "secret")
	require.NoError(t, err)
}

func TestCompareHashAndPassword_WrongPassword(t *testing.T) {
	hash, err := utils.GenerateFromPassword("secret")
	require.NoError(t, err)

	err = utils.CompareHashAndPassword(hash, "wrongpass")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid username or password")
}

// GenerateJWT & ValidateJWT

func TestGenerateJWT_And_ValidateJWT(t *testing.T) {
	token, err := utils.GenerateJWT("user123", "mysecret", time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := utils.ValidateJWT(token, "mysecret")
	require.NoError(t, err)

	require.Equal(t, "user123", claims["sub"])
}

func TestValidateJWT_InvalidSignature(t *testing.T) {
	token, err := utils.GenerateJWT("user123", "secret1", time.Hour)
	require.NoError(t, err)

	_, err = utils.ValidateJWT(token, "secret2")
	require.Error(t, err)
}

func TestValidateJWT_Expired(t *testing.T) {
	// Create immediate-expiry token
	tokenStr, err := utils.GenerateJWT("user123", "secret", time.Millisecond*1)
	require.NoError(t, err)

	time.Sleep(2 * time.Millisecond)

	_, err = utils.ValidateJWT(tokenStr, "secret")
	require.Error(t, err)
	require.Contains(t, err.Error(), "expired")
}

// Pluralize

func TestPluralize(t *testing.T) {
	require.Equal(t, "cars", utils.Pluralize("car"))
	require.Equal(t, "buses", utils.Pluralize("bus"))
	require.Equal(t, "people", utils.Pluralize("person")) // irregular plural
}

// GenerateUUID

func TestGenerateUUID(t *testing.T) {
	u1 := utils.GenerateUUID()
	u2 := utils.GenerateUUID()

	require.NotEqual(t, u1, u2)
	require.True(t, isUUID(u1))
	require.True(t, isUUID(u2))
}

func isUUID(s string) bool {
	rxUUID := regexp.MustCompile(`^[a-fA-F0-9-]{36}$`)
	return rxUUID.MatchString(s) && strings.Count(s, "-") == 4
}
