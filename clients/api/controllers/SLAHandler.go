package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-private-chaincode/clients/api/globals"
	"github.com/hyperledger/fabric-private-chaincode/clients/api/models"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils/ledger"
	t "github.com/hyperledger/fabric-private-chaincode/clients/utils/types"
	"github.com/hyperledger/fabric-private-chaincode/lib/contracts"
)

type AssetURI struct {
	ID string `uri:"id" binding:"required"`
}

func GetSingleSLA(c *gin.Context) {
	var id AssetURI
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	username := session.Get(globals.Userkey)

	user := ledger.GetUser(globals.Config, username.(string))

	if !contracts.SLAInUserContracts(user.ProviderOf, user.ClientOf, id.ID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User does not have access to this contract"})
		return
	}

	asset, err := ledger.GetSLA(globals.Config, id.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": asset})
	return
}

func GetUserSLAs(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get(globals.Userkey)

	user := ledger.GetUser(globals.Config, username.(string))

	var SLAs []t.SLA

	for _, v := range contracts.GetIDsFromString(user.ProviderOf) {
		if v == "" {
			break
		}
		asset, err := ledger.GetSLA(globals.Config, v)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		SLAs = append(SLAs, asset)
	}

	for _, v := range contracts.GetIDsFromString(user.ClientOf) {
		if v == "" {
			break
		}
		asset, err := ledger.GetSLA(globals.Config, v)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		SLAs = append(SLAs, asset)
	}
	c.JSON(http.StatusOK, gin.H{"assets": SLAs})

}

func CreateSLA(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get(globals.Userkey)

	var sla t.SLA

	if err := c.BindJSON(&sla); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sla.Provider.Name != username {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not the provider of the SLA"})
		return
	}

	err := ledger.CreateSLA(globals.Config, sla)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})

}

func GetSLAApprovalState(c *gin.Context) {
	var id AssetURI
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	username := session.Get(globals.Userkey)

	user := ledger.GetUser(globals.Config, username.(string))

	if !contracts.SLAInUserContracts(user.ProviderOf, user.ClientOf, id.ID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User does not have access to this contract"})
		return
	}
	approval, err := ledger.GetSLAApproval(globals.Config, id.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": approval})
}

func ApproveSLA(c *gin.Context) {
	var id AssetURI
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	session := sessions.Default(c)
	username := session.Get(globals.Userkey)

	user := ledger.GetUser(globals.Config, username.(string))

	if !contracts.SLAInUserContracts(user.ProviderOf, user.ClientOf, id.ID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User does not have access to this contract"})
		return
	}

	asset, err := ledger.GetSLA(globals.Config, id.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slaJSON, err := json.Marshal(asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	signature, err := utils.SignWithPrivateKey(string(slaJSON), userAPI.Mnemonic, globals.Passphrase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = ledger.Approve(globals.Config, id.ID, user.Name, signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
	return

}
