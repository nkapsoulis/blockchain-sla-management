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
	"github.com/hyperledger/fabric-private-chaincode/api/globals"
	"github.com/hyperledger/fabric-private-chaincode/api/ledger"
	"github.com/hyperledger/fabric-private-chaincode/api/middleware"
	"github.com/hyperledger/fabric-private-chaincode/api/routes"
	"github.com/hyperledger/fabric-private-chaincode/api/utils"
)

func main() {
	utils.InitConfig()
	ledger.InitLedger()

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

	router.Use(sessions.Sessions("session", store))

	public := router.Group("/")
	routes.PublicRoutes(public)

	private := router.Group("/")
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)

	router.Run(":8000")
}
