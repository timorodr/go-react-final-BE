package routes

import (

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", SignUp())
	incomingRoutes.POST("/login", Login())
}
