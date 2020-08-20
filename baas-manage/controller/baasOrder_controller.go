package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

func (a *ApiController) BaasOrderAdd(ctx *gin.Context) {
	order := new(entity.BaasOrder)
	if err := ctx.ShouldBindJSON(order); err != nil {
		gintool.ResultFail(ctx, err)
	}

	success, msg := a.baasOrderService.Add(order)
	if !success {
		gintool.ResultFail(ctx, msg)
		return
	}

	gintool.ResultOk(ctx, msg)
}

func (a *ApiController) BaasOrderUpdate(ctx *gin.Context) {
	order := new(entity.BaasOrder)
	if err := ctx.ShouldBindJSON(order); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if order.Price < 0 {
		gintool.ResultFail(ctx, common.Lang)
		return
	}

	success, msg := a.baasOrderService.Update(order)
	if !success {
		gintool.ResultFailData(ctx, msg, nil)
		return
	}

	gintool.ResultOk(ctx, msg)
}

func (a *ApiController) BaasOrderDeleted(ctx *gin.Context) {
	order := new(entity.BaasOrder)
	if err := ctx.ShouldBindJSON(order); err != nil {
		gintool.ResultFail(ctx, err)
	}

	success, msg := a.baasOrderService.Delete(order.Id)
	if !success {
		gintool.ResultFail(ctx, msg)
		return
	}

	gintool.ResultOk(ctx, msg)
}

func (a *ApiController) BaasOrderGetList(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "page_error", nil, common.Lang))
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "limit_error", nil, common.Lang))
		return
	}

	status := ctx.Query("status")
	tagName := ctx.Query("tagName")

	has, baasOrders, total := a.baasOrderService.GetList(status, &entity.BaasOrder{TagName: tagName}, page, limit)
	if !has {
		gintool.ResultFail(ctx, ii18n.T("app", "get_order_fail", nil, common.Lang))
		return
	}
	gintool.ResultList(ctx, baasOrders, total)
}
