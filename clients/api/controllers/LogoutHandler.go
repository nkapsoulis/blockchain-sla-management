package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-private-chaincode/clients/api/globals"
)

func PostLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(globals.Userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	return
}
