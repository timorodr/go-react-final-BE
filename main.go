package main

import (
	"os"


	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "github.com/sashabaranov/go-openai" // import our own routes can be internal or external
	middleware "github.com/timorodr/go-react-final/server/middleware"
	"github.com/timorodr/go-react-final/server/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	// c := openai.NewClient(os.Getenv("OPENAI_KEY"))
	// config := cors.DefaultConfig()
	config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"} // Allow requests from localhost:3000
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // Allow specified methods
    config.AllowHeaders = []string{"Authorization", "Content-Type", "Token"} // Allow Authorization header

    router.Use(cors.New(config))
	// config.AllowOrigins = []string{"http://localhost:3000"}
	// config.AllowAllOrigins = true

	router.Use(gin.Logger()) // shows when whcih API was called
	// router.Use(cors.Default())
	// router.Use(cors.New(config))
	// router.Use(cors.Default())
	routes.UserRoutes(router)

	router.Use(middleware.Authentication())

	
	// authorized := router.Group("/user")
	// authorized.Use(middleware.Authentication())

	// authorized.GET("/entries", routes.GetEntries) //

	router.POST("/user/entry/create/:id", routes.AddEntry())
	router.GET("/user/entries/:id", routes.GetEntries)
	router.POST("/user/logout", routes.Logout)
	router.PUT("/user/entry/update/:id/:medication_id", routes.UpdateEntry)
	router.DELETE("/user/entry/delete/:id/:medication_id", routes.DeleteEntry)

	// router.POST("/entry/create", routes.AddEntry)
	// authorized.GET("/entry/:id/", routes.GetEntryById)

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
