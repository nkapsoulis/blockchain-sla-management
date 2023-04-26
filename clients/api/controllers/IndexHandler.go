package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-private-chaincode/clients/api/globals"
)

// GetName returns the name of this API
func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"name": globals.AppName})
	return
}
