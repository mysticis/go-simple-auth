package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//JWTWrapper wraps the signing key and issuer
type JWTWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

//JWTClaim adds email as claim to the token

type JWTClaim struct {
	Email string
	jwt.StandardClaims
}

//GenerateToken generates a jwt token

func (j *JWTWrapper) GenerateToken(email string) (signedToken string, err error) {

	claims := &JWTClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))

	if err != nil {
		return
	}

	return
}

//VaidateToken validates the generatedtoken

func (j *JWTWrapper) ValidateToken(signedToken string) (claims *JWTClaim, err error) {

	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)

	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired!")
		return
	}

	return
}
