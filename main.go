package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/sashabaranov/go-openai" // import our own routes can be internal or external
	"github.com/timorodr/go-react-final/server/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// c := openai.NewClient(os.Getenv("OPENAI_KEY"))
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"https://go-react-fe.netlify.app"}

	router := gin.New()
	router.Use(gin.Logger()) // shows when whcih API was called
	router.Use(cors.Default())
	// router.Use(cors.New(config))

	routes.UserRoutes(router)

	router.POST("/entry/create", routes.AddEntry)
	router.GET("/entries", routes.GetEntries)
	router.GET("/entry/:id/", routes.GetEntryById)
	// router.GET("/ingredient/:ingredient", routes.GetEntriesByIngredient)

	router.PUT("/entry/update/:id", routes.UpdateEntry)
	// router.PUT("/ingredient/update/:id", routes.UpdateIngredient)
	router.DELETE("/entry/delete/:id", routes.DeleteEntry)

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})


	router.Run(":" + port)
}
