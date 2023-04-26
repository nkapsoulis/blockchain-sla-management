package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-private-chaincode/clients/api/controllers"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/", controllers.GetIndex)
	g.POST("/auth/login", controllers.Login)
}

func PrivateRoutes(g *gin.RouterGroup) {
	// Auth routes
	g.GET("/auth/user", controllers.GetUserData)
	g.POST("/auth/logout", controllers.PostLogout)

	// Asset routes
	g.GET("/assets", controllers.GetUserSLAs)
	g.POST("/assets", controllers.CreateSLA)

	g.GET("/assets/:id", controllers.GetSingleSLA)
	g.GET("/assets/:id/approvals", controllers.GetSLAApprovalState)
	g.POST("/assets/:id/approve", controllers.ApproveSLA)
}
