package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mysticis/go-net/database"
	"github.com/mysticis/go-net/models"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {

	var actualResult models.User

	user := models.User{
		Name:     "martin",
		Email:    "martin@test.com",
		Password: "secret",
	}

	payload, err := json.Marshal(&user)

	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/api/public/signup", bytes.NewBuffer(payload))

	assert.NoError(t, err)

	responseRecorder := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(responseRecorder)

	ctx.Request = request

	err = database.InitDatabase()

	assert.NoError(t, err)

	database.GlobalDB.AutoMigrate(&models.User{})

	SignUp(ctx)

	assert.Equal(t, 200, responseRecorder.Code)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &actualResult)

	assert.NoError(t, err)

	assert.Equal(t, user.Name, actualResult.Name)

	assert.Equal(t, user.Email, actualResult.Email)

}

func TestSignUpInvalidJSON(t *testing.T) {

	user := "test"

	payload, err := json.Marshal(&user)

	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/public/signup", bytes.NewBuffer(payload))

	assert.NoError(t, err)

	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = req

	SignUp(ctx)

	assert.Equal(t, 400, w.Code)

}

func TestLogin(t *testing.T) {

	user := LoginPayload{
		Email:    "martin@test.com",
		Password: "secret",
	}

	payload, err := json.Marshal(&user)

	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/public/login", bytes.NewBuffer(payload))

	assert.NoError(t, err)

	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = req

	err = database.InitDatabase()

	assert.NoError(t, err)

	err = database.GlobalDB.AutoMigrate(&models.User{})

	assert.NoError(t, err)

	Login(ctx)

	assert.Equal(t, 200, w.Code)

}

func TestLoginInvalidJSON(t *testing.T) {

	user := "test"

	payload, err := json.Marshal(&user)

	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/public/login", bytes.NewBuffer(payload))

	assert.NoError(t, err)

	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = req

	Login(ctx)

	assert.Equal(t, 400, w.Code)
}

func TestInvalidLoginCredentials(t *testing.T) {

	user := LoginPayload{
		Email:    "martin@test.com",
		Password: "invalid",
	}

	payload, err := json.Marshal(&user)

	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/public/login", bytes.NewBuffer(payload))

	assert.NoError(t, err)

	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = req

	err = database.InitDatabase()

	assert.NoError(t, err)

	err = database.GlobalDB.AutoMigrate(&models.User{})

	assert.NoError(t, err)

	Login(ctx)

	assert.Equal(t, 401, w.Code)

	database.GlobalDB.Unscoped().Where("email = ?", user.Email).Delete(&models.User{})
}
