package controller

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/config"
	"github.com/syyongx/ii18n"

	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-gateway/service"
)

type ApiController struct {
	chainService     *service.ChainService
	channelService   *service.ChannelService
	chaincodeService *service.ChaincodeService
	dashboardService *service.DashboardService
	userService      *service.UserService
	roleService      *service.RoleService
	baasComboService *service.BaasComboService
	baasOrderService *service.BaasOrderService
}

func NewApiController(userService *service.UserService, roleService *service.RoleService, chainService *service.ChainService, channelService *service.ChannelService, chaincodeService *service.ChaincodeService, baasComboService *service.BaasComboService, baasOrderService *service.BaasOrderService, dashboardService *service.DashboardService) *ApiController {
	return &ApiController{
		userService:      userService,
		roleService:      roleService,
		chainService:     chainService,
		channelService:   channelService,
		chaincodeService: chaincodeService,
		baasComboService: baasComboService,
		baasOrderService: baasOrderService,
		dashboardService: dashboardService,
	}
}

func (a *ApiController) Upload(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	pointCount := strings.Count(file.Filename, ".")
	if pointCount > 1 || !strings.HasSuffix(file.Filename, ".go") {
		gintool.ResultFail(ctx, ii18n.T("app", "fileTypeErr", nil, common.Lang))
		return
	}
	if file.Size > 1024*50 {
		gintool.ResultFail(ctx, ii18n.T("app", "fileSizeErr", nil, common.Lang))
		return
	}
	chaincodePath := config.Config.GetString("ChaincodePath")

	exist, err := common.PathExists(chaincodePath)
	if !exist {
		err = os.Mkdir(chaincodePath, os.ModePerm)
		if err != nil {
			gintool.ResultFail(ctx, ii18n.T("app", "mkdirFail", nil, common.Lang))
			return
		}
	}

	path := fmt.Sprintf("%s%d", chaincodePath, time.Now().UnixNano())
	ctx.SaveUploadedFile(file, path)
	gintool.ResultOk(ctx, path)

}
