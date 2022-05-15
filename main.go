package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mysticis/go-net/controllers"
	"github.com/mysticis/go-net/database"
	"github.com/mysticis/go-net/middleware"
	"github.com/mysticis/go-net/models"
)

func setUpRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {

		c.String(200, "pong")
	})

	//group api

	api := r.Group("/api")

	{
		public := api.Group("/public")

		{
			public.POST("/login", controllers.Login)
			public.POST("/signup", controllers.SignUp)

		}

		//protected rotes go here
		protected := api.Group("/protected").Use(middleware.Authz())

		{
			protected.GET("/profile", controllers.Profile)
		}
	}
	return r
}

func main() {

	err := database.InitDatabase()

	if err != nil {
		log.Fatalln("could not initialise database", err)
	}

	database.GlobalDB.AutoMigrate(&models.User{})

	r := setUpRouter()

	r.Run(":8080")
}
