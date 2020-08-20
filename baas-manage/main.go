package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-core/common/xorm"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/config"
	"github.com/paybf/baasmanager/baas-gateway/controller"
	"github.com/paybf/baasmanager/baas-gateway/service"
)

func main() {
	dbengine := xorm.GetEngine(config.Config.GetString("BaasGatewayDbconfig"))
	common.ConnRedis()
	defer common.RedisClient.Close()
	fabricService := service.NewFabricService()
	userService := service.NewUserService(dbengine)
	apiController := controller.NewApiController(
		service.NewUserService(dbengine),
		service.NewRoleService(dbengine),
		service.NewChainService(dbengine, fabricService, userService),
		service.NewChannelService(dbengine, fabricService),
		service.NewChaincodeService(dbengine, fabricService),
		service.NewComboService(dbengine),
		service.NewOrderService(dbengine),
		service.NewDashboardService(dbengine),
	)
	router := gin.New()
	router.Use(gintool.Logger())
	router.Use(gin.Recovery())
	router.Use(common.Cors())
	gintool.UseSession(router)
	router.Use(common.I18nConfig())

	api := router.Group("/api")
	{
		api.POST("/user/login", apiController.UserLogin)
		api.POST("/user/logout", apiController.UserLogout)

		api.Use(apiController.UserAuthorize)

		api.GET("/user/info", apiController.UserInfo)

		api.GET("/user/list", apiController.UserList)
		api.POST("/user/add", apiController.UserAdd)
		api.POST("/user/updateAuth", apiController.UserUpdateAuth)
		api.POST("/user/addAuth", apiController.UserAddAuth)
		api.POST("/user/delAuth", apiController.UserDelAuth)
		api.POST("/user/update", apiController.UserUpdate)
		api.POST("/user/delete", apiController.UserDelete)
		api.GET("/user/useridList", apiController.GetUserList)
		api.POST("/user/changePwd", apiController.ChangePwd)
		api.POST("/user/resetPwd", apiController.ResetPwd)

		api.GET("/role/list", apiController.RoleList)
		api.GET("/role/allList", apiController.RoleAllList)
		api.POST("/role/add", apiController.RoleAdd)
		api.POST("/role/update", apiController.RoleUpdate)
		api.POST("/role/delete", apiController.RoleDelete)

		api.GET("/chain/list", apiController.ChainList)
		api.POST("/chain/add", apiController.ChainAdd)
		api.POST("/chain/update", apiController.ChainUpdate)
		api.POST("/chain/get", apiController.ChainGet)
		api.POST("/chain/delete", apiController.ChainDeleted)
		api.POST("/chain/build", apiController.ChainBuild)
		api.POST("/chain/run", apiController.ChainRun)
		api.POST("/chain/stop", apiController.ChainStop)
		api.POST("/chain/release", apiController.ChainRelease)
		api.POST("/chain/changeSize", apiController.ChangeChainResouces)
		api.GET("/chain/download", apiController.ChainDownload)
		api.GET("/chain/podsQuery", apiController.ChainPodsQuery)

		api.POST("/channel/add", apiController.ChannelAdd)
		api.POST("/channel/get", apiController.ChannelGet)
		api.GET("/channel/allList", apiController.ChannelAll)

		api.GET("/chaincode/list", apiController.ChaincodeList)
		api.POST("/chaincode/add", apiController.ChaincodeAdd)
		api.POST("/chaincode/deploy", apiController.ChaincodeDeploy)
		api.POST("/chaincode/upgrade", apiController.ChaincodeUpgrade)
		api.POST("/chaincode/query", apiController.ChaincodeQuery)
		api.GET("/chaincode/queryLedger", apiController.ChaincodeLedgerQuery)
		api.GET("/chaincode/queryLatestBlocks", apiController.ChaincodeLatestBlocksQuery)
		api.GET("/chaincode/queryBlock", apiController.ChaincodeBlockQuery)
		api.POST("/chaincode/invoke", apiController.ChaincodeInvoke)
		api.POST("/chaincode/get", apiController.ChaincodeGet)
		api.POST("/chaincode/delete", apiController.ChaincodeDeleted)

		api.GET("/combo/get", apiController.BaasComboGetByDisplayName)
		api.GET("/combo/list", apiController.BaasComboGetAllList)
		api.POST("/combo/add", apiController.BaasComboAdd)
		api.POST("/combo/update", apiController.BaasComboUpdate)
		api.POST("/combo/delete", apiController.BaasComboDelete)

		api.GET("/order/list", apiController.BaasOrderGetList)
		api.POST("/order/add", apiController.BaasOrderAdd)
		api.POST("/order/update", apiController.BaasOrderUpdate)
		api.POST("/order/del", apiController.BaasOrderDeleted)

		api.POST("/upload", apiController.Upload)

		api.GET("/dashboard/counts", apiController.DashboardCounts)
		api.GET("/dashboard/sevenDays", apiController.DashboardSevenDays)
		api.GET("/dashboard/consensusTotal", apiController.DashboardConsensusTotal)
	}

	router.Run(":" + config.Config.GetString("BaasGatewayPort"))
}
