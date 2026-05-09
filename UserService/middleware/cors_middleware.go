package middleware

import "github.com/gin-gonic/gin"

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		   		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		   		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		   		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		   		if c.Request.Method == "OPTIONS" {
		   			c.AbortWithStatus(204)
		   			return
		   		}

		   		c.Next() */
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Replace * with your domain for security
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
