package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/syyongx/ii18n"
)

func (a *ApiController) DashboardCounts(ctx *gin.Context) {
	isSuccess, ash := a.dashboardService.Counts()
	if isSuccess {
		gintool.ResultOk(ctx, ash)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "stat_usersChains_chaincodes_channels_error", nil, common.Lang))
	}
}

func (a *ApiController) DashboardConsensusTotal(ctx *gin.Context) {
	isSuccess, ash := a.dashboardService.ConsensusTotal()
	if isSuccess {
		gintool.ResultOk(ctx, ash)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "stat_consensusTotal_error", nil, common.Lang))
	}
}

func (a *ApiController) DashboardSevenDays(ctx *gin.Context) {
	start, err := strconv.Atoi(ctx.Query("start"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "start_error", nil, common.Lang))
		return
	}

	end, err := strconv.Atoi(ctx.Query("end"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "end_error", nil, common.Lang))
		return
	}

	isSuccess, ash := a.dashboardService.SevenDays(start, end)
	if isSuccess {
		gintool.ResultOk(ctx, ash)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "stat_seven_dat_data_error", nil, common.Lang))
	}
}
