package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mysticis/go-net/auth"
	"github.com/mysticis/go-net/controllers"
	"github.com/mysticis/go-net/database"
	"github.com/mysticis/go-net/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthzNoHeader(t *testing.T) {

	router := gin.Default()

	router.Use(Authz())

	router.GET("/api/protected/profile", controllers.Profile)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/protected/profile", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}

func TestAuthzInvalidTokenFormat(t *testing.T) {

	router := gin.Default()

	router.Use(Authz())

	router.GET("/api/protected/profile", controllers.Profile)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/protected/profile", nil)

	req.Header.Add("Authorization", "test")

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

}

func TestAuthzInvalidToken(t *testing.T) {

	invalidToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	router := gin.Default()

	router.Use(Authz())

	router.GET("/api/protected/profile", controllers.Profile)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/protected/profile", nil)

	req.Header.Add("Authorization", invalidToken)

	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

}

func TestValidToken(t *testing.T) {

	var response models.User

	err := database.InitDatabase()

	assert.NoError(t, err)

	err = database.GlobalDB.AutoMigrate(&models.User{})

	assert.NoError(t, err)

	user := models.User{
		Name:     "xavi",
		Email:    "xavi@test.org",
		Password: "secret",
	}

	jwtWrapper := auth.JWTWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}
	token, err := jwtWrapper.GenerateToken(user.Email)

	assert.NoError(t, err)

	err = user.HashPassword(user.Password)

	assert.NoError(t, err)

	result := database.GlobalDB.Create(&user)

	assert.NoError(t, result.Error)

	router := gin.Default()

	router.Use(Authz())

	router.GET("/api/protected/profile", controllers.Profile)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/protected/profile", nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	router.ServeHTTP(w, req)

	err = json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)

	assert.Equal(t, "xavi@test.org", response.Email)

	assert.Equal(t, "xavi", response.Name)

	database.GlobalDB.Unscoped().Where("email = ?", user.Email).Delete(&models.User{})
}
