package routes

import (
	// controller "github.com/timorodr/go-react-final/server/controllers"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", SignUp())
	incomingRoutes.POST("/login", Login())
}
