package middlewares

import (
	"log"
	"os"

	"github.com/Habeebamoo/Clivo/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthenticateAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization") 
		if token == "" {
			utils.Abort(c, 401, "Token Missing", nil)
			return 
		}

		log.Println(token)

		//verify token
		if token != os.Getenv("AGENT_TOKEN") {
			utils.Abort(c, 401, "Invalid Token", nil)
			return 
		}

		//call the next middleware
		c.Next()
	}
}