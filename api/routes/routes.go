package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-private-chaincode/api/controllers"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/", controllers.GetIndex)
	g.POST("/login", controllers.Login)
}

func PrivateRoutes(g *gin.RouterGroup) {
	// Auth routes
	g.GET("/user", controllers.GetUserData)
	g.POST("/logout", controllers.PostLogout)

	// Asset routes
	g.GET("/assets/:id", controllers.GetSingleSLA)
	g.POST("/assets", controllers.CreateSLA)
	g.GET("/assets/:id/approvals", controllers.GetSLAApprovalState)
	g.POST("/assets/:id/approve", controllers.ApproveSLA)
}
