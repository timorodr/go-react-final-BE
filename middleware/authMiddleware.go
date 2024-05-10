package middleware

import (
	"fmt"
	"net/http"
	"strings"

	helper "github.com/timorodr/go-react-final/server/helpers"

	"github.com/gin-gonic/gin"
)

// Authz validates token and authorizes users
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.Request.Header)
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

        ClientToken := authParts[1]

		claims, err := helper.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		// c.Set("user_id", claims.User_id)
		// c.Set("first_name", claims.First_name)
		// c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)

		c.Next()

	}
}

// func NoCache() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         c.Header("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
//         c.Header("Pragma", "no-cache") // HTTP 1.0.
//         c.Header("Expires", "0") // Proxies.
//         c.Next()
//     }
// }

func AuthenticationToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println(c.Request.Header)
        clientToken := c.Request.Header.Get("Token")
        if clientToken == "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
            c.Abort()
            return
        }

        claims, err := helper.ValidateToken(clientToken)
        if err != "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err})
            c.Abort()
            return
        }

        c.Set("email", claims.Email)
        // c.Set("user_id", claims.User_id)
        // c.Set("first_name", claims.First_name)
        // c.Set("last_name", claims.Last_name)
        c.Set("user_id", claims.Uid)

        c.Next()

    }
}