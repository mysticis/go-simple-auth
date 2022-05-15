package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {

	jwtWrapper := JWTWrapper{
		SecretKey:       "verysecretKey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	generatedToken, err := jwtWrapper.GenerateToken("joe@test.org")

	assert.NoError(t, err)

	os.Setenv("testtoken", generatedToken)
}

func TestVaidateToken(t *testing.T) {

	encodedToken := os.Getenv("testtoken")

	jwtWrapper := JWTWrapper{
		SecretKey: "verysecretKey",
		Issuer:    "AuthService",
	}

	claims, err := jwtWrapper.ValidateToken(encodedToken)

	assert.NoError(t, err)

	assert.Equal(t, "joe@test.org", claims.Email)

	assert.Equal(t, "AuthService", claims.Issuer)
}
