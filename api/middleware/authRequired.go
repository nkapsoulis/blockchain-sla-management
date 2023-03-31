package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"log"
	"net/http"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		log.Println("User not logged in")
		c.JSON(http.StatusForbidden, gin.H{"error": "Not Authenticated"})
		c.Abort()
		return
	}
	c.Next()
}
