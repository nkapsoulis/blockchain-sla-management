/*
Copyright IBM Corp. All Rights Reserved.
Copyright 2020 Intel Corporation

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-private-chaincode/clients/api/globals"
	"github.com/hyperledger/fabric-private-chaincode/clients/api/middleware"
	"github.com/hyperledger/fabric-private-chaincode/clients/api/routes"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils/ledger"
)

func main() {
	globals.Config = utils.InitConfig()
	ledger.InitLedger(globals.Config, globals.Passphrase)

	router := gin.Default()

	// Configure the cookies used
	store := cookie.NewStore(globals.Secret)
	store.Options(sessions.Options{
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   false,
		MaxAge:   60 * 60 * 24 * 30, // one month in seconds
		Path:     "/",
		Domain:   "localhost",
	})

	router.Use(sessions.Sessions(globals.SessionName, store))

	public := router.Group("/")
	routes.PublicRoutes(public)

	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)

	router.Run(":8000")
}
