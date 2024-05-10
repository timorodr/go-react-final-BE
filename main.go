package main

import (
	"os"


	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	middleware "github.com/timorodr/go-react-final/server/middleware"
	"github.com/timorodr/go-react-final/server/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	
	config := cors.DefaultConfig()
    config.AllowOrigins = []string{"https://main--hilarious-biscotti-0d1872.netlify.app/"} // Allow requests from localhost:3000
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // Allow specified methods
    config.AllowHeaders = []string{"Authorization", "Content-Type", "Token"} // Allow Authorization header

    router.Use(cors.New(config))

	router.Use(gin.Logger()) // shows when whcih API was called
	routes.UserRoutes(router)

	router.Use(middleware.Authentication())

	

	router.POST("/user/entry/create/:id", routes.AddEntry(), middleware.Authentication())
	router.GET("/user/entries/:id", routes.GetEntries, middleware.Authentication())
	router.POST("/user/logout", routes.Logout)
	router.PUT("/user/entry/update/:id/:medication_id", routes.UpdateEntry, middleware.Authentication())
	router.DELETE("/user/entry/delete/:id/:medication_id", routes.DeleteEntry, middleware.Authentication())


	router.Run(":" + port)
}
