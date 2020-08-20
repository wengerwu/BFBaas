package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
	tsgutils "github.com/typa01/go-utils"
)

func (a *ApiController) BaasComboAdd(ctx *gin.Context) {
	baasCombo := new(entity.BaasCombo)
	if err := ctx.ShouldBindJSON(baasCombo); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if baasCombo.DisplayName == "" || baasCombo.Remark == "" || baasCombo.Price < 0 || baasCombo.Sort < 0 {
		gintool.ResultFail(ctx, common.Lang)
		return
	}

	baasCombo.Id = tsgutils.GUID()
	isSuccess, msg := a.baasComboService.Add(baasCombo)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) BaasComboUpdate(ctx *gin.Context) {
	baasCombo := new(entity.BaasCombo)
	if err := ctx.ShouldBindJSON(baasCombo); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if baasCombo.DisplayName == "" || baasCombo.Remark == "" || baasCombo.Price < 0 || baasCombo.Sort < 0 {
		gintool.ResultFail(ctx, common.Lang)
		return
	}

	isSuccess, msg := a.baasComboService.Update(baasCombo)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFailData(ctx, msg, nil)
	}
}

func (a *ApiController) BaasComboDelete(ctx *gin.Context) {
	baasCombo := new(entity.BaasCombo)
	if err := ctx.ShouldBindJSON(baasCombo); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.baasComboService.Delete(baasCombo.Id)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) BaasComboGetAllList(ctx *gin.Context) {
	isSuccess, list := a.baasComboService.GetAllList()
	if isSuccess {
		gintool.ResultList(ctx, list, 0)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "get_combo_fail", nil, common.Lang))
	}
}

func (a *ApiController) BaasComboGetByDisplayName(ctx *gin.Context) {
	displayName := ctx.Query("displayName")
	if displayName == "" {
		a.BaasComboGetAllList(ctx)
		return
	}
	isSuccess, combo := a.baasComboService.GetCombo(displayName)
	if isSuccess {
		gintool.ResultOk(ctx, combo)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "get_combo_fail", nil, common.Lang))
	}
}
