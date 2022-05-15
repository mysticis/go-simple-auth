package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mysticis/go-net/database"
	"github.com/mysticis/go-net/models"
	"gorm.io/gorm"
)

//Profile returns user data

func Profile(c *gin.Context) {

	var user models.User

	email, _ := c.Get("email") //from the authorization middleware

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)

	if result.Error == gorm.ErrRecordNotFound {

		c.JSON(404, gin.H{"msg": "user not found"})

		c.Abort()

		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{"msg": "could not get user profile"})
		c.Abort()
		return
	}

	user.Password = " "

	c.JSON(200, user)

	return
}
