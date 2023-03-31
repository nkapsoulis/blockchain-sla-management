package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-private-chaincode/api/globals"
	"github.com/hyperledger/fabric-private-chaincode/api/ledger"
)

func GetUserData(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get(globals.Userkey)

	user := ledger.GetUser(username.(string))

	c.JSON(http.StatusOK, gin.H{"user": user})
}
