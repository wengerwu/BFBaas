package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

func (a *ApiController) RoleList(ctx *gin.Context) {
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

	name := ctx.Query("name")
	b, list, total := a.roleService.GetList(&entity.Role{Name: name}, page, limit)
	if b {
		gintool.ResultList(ctx, list, total)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "get_page_role_list_error", nil, common.Lang))
	}
}

func (a *ApiController) RoleAllList(ctx *gin.Context) {
	b, list := a.roleService.GetAll()
	if b {
		gintool.ResultOk(ctx, list)

	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "get_all_role_list_error", nil, common.Lang))
	}
}

func (a *ApiController) RoleAdd(ctx *gin.Context) {
	role := new(entity.Role)
	if err := ctx.ShouldBindJSON(role); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if role.Name == "" || role.Rkey == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "notValue", nil, common.Lang))
		return
	}

	isSuccess, msg := a.roleService.Add(role)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) RoleUpdate(ctx *gin.Context) {
	role := new(entity.Role)
	if err := ctx.ShouldBindJSON(role); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if role.Name == "" || role.Rkey == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "notValue", nil, common.Lang))
		return
	}

	isSuccess, msg := a.roleService.Update(role)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) RoleDelete(ctx *gin.Context) {
	role := new(entity.Role)
	if err := ctx.ShouldBindJSON(role); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if role.Rkey == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "notValue", nil, common.Lang))
		return
	}

	isSuccess, msg := a.roleService.Delete(role.Rkey)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}
