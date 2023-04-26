package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-private-chaincode/clients/api/globals"
	"github.com/hyperledger/fabric-private-chaincode/clients/api/models"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils/ledger"
	"github.com/hyperledger/fabric-private-chaincode/lib"
)

// login ensures the user is logged in
func Login(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get(globals.Userkey)
	if user != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Already authenticated"})
		return
	}

	var userAPI models.User

	if err := c.BindJSON(&userAPI); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	if userAPI.Mnemonic == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mnemonic is missing"})
		return
	}

	if userAPI.Name == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Name is missing"})
		return
	}

	userLedger := ledger.GetUser(globals.Config, userAPI.Name)

	if userLedger == (lib.User{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	match, err := utils.PubKeyMatchesMnemonic(userAPI.Mnemonic, globals.Passphrase, userLedger.PubKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !match {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mnemonic and public key do not match"})
		return
	}

	session.Set(globals.Userkey, userAPI.Name)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": userAPI.Name})
	return
}
