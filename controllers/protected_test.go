package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mysticis/go-net/database"
	"github.com/mysticis/go-net/models"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {

	var profile models.User

	err := database.InitDatabase()

	assert.NoError(t, err)

	database.GlobalDB.AutoMigrate(&models.User{})

	user := models.User{
		Email:    "xavier@test.com",
		Name:     "xavier",
		Password: "secret",
	}

	err = user.HashPassword(user.Password)

	assert.NoError(t, err)

	err = user.CreateUserRecord()

	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/api/protected/profile", nil)

	assert.NoError(t, err)

	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = req

	ctx.Set("email", "xavier@test.com")

	Profile(ctx)

	err = json.Unmarshal(w.Body.Bytes(), &profile)

	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)

	log.Println(profile)

	assert.Equal(t, user.Email, profile.Email)

	assert.Equal(t, user.Name, profile.Name)
}

func TestProfileNotFound(t *testing.T) {

	var profile models.User

	err := database.InitDatabase()

	assert.NoError(t, err)

	database.GlobalDB.AutoMigrate(&models.User{})

	req, err := http.NewRequest("GET", "/api/protected/profile", nil)

	assert.NoError(t, err)

	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = req

	ctx.Set("email", "notfound@gmail.com")

	Profile(ctx)

	err = json.Unmarshal(w.Body.Bytes(), &profile)

	assert.NoError(t, err)

	assert.Equal(t, 404, w.Code)

	database.GlobalDB.Unscoped().Where("email = ?", "xavier@test.com").Delete(&models.User{})

}
