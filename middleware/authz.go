package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mysticis/go-net/auth"
)

func Authz() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("Authorization")

		if clientToken == "" {
			ctx.JSON(403, gin.H{"msg": "no authorization header provided"})
			ctx.Abort()
			return
		}
		extractedToken := strings.Split(clientToken, "Bearer ")

		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			ctx.JSON(400, "incorrect format of authorization token")
			ctx.Abort()
			return
		}

		jwtWrapper := auth.JWTWrapper{
			SecretKey: "verysecretkey",
			Issuer:    "AuthService",
		}

		claims, err := jwtWrapper.ValidateToken(clientToken)

		if err != nil {
			ctx.JSON(401, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)

		ctx.Next()
	}
}
