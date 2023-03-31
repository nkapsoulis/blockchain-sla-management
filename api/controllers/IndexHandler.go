package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-private-chaincode/api/globals"
	"github.com/hyperledger/fabric-private-chaincode/api/models"
)

// GetName returns the name of this API
func GetIndex(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, models.Response{Name: globals.AppName})
	return
}
