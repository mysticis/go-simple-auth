package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mysticis/go-net/auth"
	"github.com/mysticis/go-net/database"
	"github.com/mysticis/go-net/models"
	"gorm.io/gorm"
)

func SignUp(c *gin.Context) {

	var user models.User

	//get details from the user
	err := c.ShouldBindJSON(&user)

	if err != nil {

		log.Println(err)

		c.JSON(400, gin.H{"msg": "invalid json"})

		c.Abort()
		return
	}

	// if all goes well, try to hash the password

	err = user.HashPassword(user.Password)

	if err != nil {
		log.Println(err.Error())

		c.JSON(500, gin.H{"msg": "error hasshing password"})

		c.Abort()

		return
	}

	//if hashing succeeds, we store the user in the database

	err = user.CreateUserRecord()

	if err != nil {

		log.Println(err)

		c.JSON(500, gin.H{"msg": "could not create user"})

		c.Abort()

		return
	}

	//if all goes well we store the user and return the stores user to the client

	c.JSON(200, user)

}

//after signup, the user will login with email and username

//Login Payload body

type LoginPayload struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

//LoginResponse token sent from server

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(c *gin.Context) {

	var user models.User

	var payload LoginPayload

	//get details from request body
	err := c.ShouldBindJSON(&payload)

	if err != nil {
		c.JSON(400, gin.H{"msg": "invalid json"})
		c.Abort()
		return
	}

	//if all goes well get the stored user from the database

	dBResult := database.GlobalDB.Where("email = ?", payload.Email).First(&user)

	if dBResult.Error == gorm.ErrRecordNotFound {

		c.JSON(401, gin.H{"msg": "user not found or invalid credentials"})

		c.Abort()

		return
	}

	//if user is found, we check that the supplied password matches the stored hashedpassword

	err = user.CheckPassword(payload.Password)

	if err != nil {
		log.Println(err)
		c.JSON(401, gin.H{"msg": "invalid user credentials"})

		c.Abort()
		return
	}

	jwtWrapper := auth.JWTWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	signedToken, err := jwtWrapper.GenerateToken(payload.Email)

	if err != nil {

		log.Println(err)

		c.JSON(500, gin.H{"msg": "error signing token"})

		c.Abort()

		return
	}

	//if token signing succeeds, we send it to the client

	tokenResponse := LoginResponse{
		Token: signedToken,
	}

	c.JSON(200, tokenResponse)

	return

}
