package internal

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowHeadersCORS = []string{
	"Content-Type", "Content-Length", "Accept-Encoding", "Authorization",
	"X-Header-TraxApp",
}

var allowMethodsCORS = []string{"OPTIONS", "GET", "POST", "PATCH", "PUT"}

func SetCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(allowHeadersCORS, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(allowMethodsCORS, ", "))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
