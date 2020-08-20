package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-core/core/model"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

func (a *ApiController) ChainAdd(ctx *gin.Context) {
	chain := new(entity.Chain)
	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	nameRule, _ := regexp.MatchString("^[A-Za-z0-9]+$", chain.Name)
	if !nameRule || chain.Name == "" || chain.Consensus == "" || chain.OrderCount <= 0 || chain.PeerCount <= 0 || chain.PeersOrgs == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "args_fail", nil, common.Lang))
		return
	}

	//验证组织是否有重复，不符合规范的数据
	peerOrgsArry := strings.Split(chain.PeersOrgs, ",")
	newArr := make([]string, 0)
	for i := 0; i < len(peerOrgsArry); i++ {
		text := peerOrgsArry[i]
		textRules, _ := regexp.MatchString("^[A-Za-z0-9]+$", text)
		if !textRules {
			gintool.ResultFail(ctx, ii18n.T("app", "args_fail", nil, common.Lang))
			return
		}
		if len(newArr) < 1 {
			newArr = append(newArr, text)
			continue
		}
		for j := 0; j < len(newArr); j++ {
			newArrText := newArr[j]
			if text == newArrText {
				gintool.ResultFail(ctx, ii18n.T("app", "args_fail", nil, common.Lang))
				return
			}
		}
		newArr = append(newArr, text)
	}

	isHaveName := a.chainService.IsHaveName(chain, true)
	if isHaveName {
		gintool.ResultFail(ctx, ii18n.T("app", "chainNameIsHave", nil, common.Lang))
		return
	}

	chain.Created = time.Now().Unix()
	isSuccess, msg := a.chainService.Add(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainGet(ctx *gin.Context) {
	chain := new(entity.Chain)
	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, chain := a.chainService.GetByChain(chain)
	if isSuccess {
		gintool.ResultOk(ctx, chain)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "get_chain_fail", nil, common.Lang))
	}
}

func (a *ApiController) ChainList(ctx *gin.Context) {
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

	b, list, total := a.chainService.GetList(&entity.Chain{
		Name: name,
	}, page, limit)

	if b {
		gintool.ResultList(ctx, list, total)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "get_chain_list_fail", nil, common.Lang))
	}
}

func (a *ApiController) GetUserList(ctx *gin.Context) {
	userNameOrPhone := ctx.Query("userNameOrPhone")
	useridList := a.chainService.GetUserIdList(userNameOrPhone)
	gintool.ResultOk(ctx, useridList)
}

func (a *ApiController) ChainUpdate(ctx *gin.Context) {
	chain := new(entity.Chain)
	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	nameRule, _ := regexp.MatchString("^[A-Za-z0-9]+$", chain.Name)
	if !nameRule || chain.Name == "" || chain.Consensus == "" || chain.OrderCount <= 0 || chain.PeerCount <= 0 || chain.PeersOrgs == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "args_fail", nil, common.Lang))
		return
	}

	isHaveName := a.chainService.IsHaveName(chain, false)
	if isHaveName {
		gintool.ResultFail(ctx, ii18n.T("app", "chainNameIsHave", nil, common.Lang))
		return
	}

	isSuccess, msg := a.chainService.Update(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainDeleted(ctx *gin.Context) {
	chain := new(entity.Chain)
	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chainService.Delete(chain.Id)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainBuild(ctx *gin.Context) {
	chain := new(entity.Chain)
	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chainService.BuildChain(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}

}

func (a *ApiController) ChainRun(ctx *gin.Context) {
	chain := new(entity.Chain)
	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chainService.RunChain(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}

}

func (a *ApiController) ChainStop(ctx *gin.Context) {
	chain := new(entity.Chain)
	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chainService.StopChain(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainRelease(ctx *gin.Context) {
	chain := new(entity.Chain)
	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chainService.ReleaseChain(chain)
	if isSuccess {
		a.channelService.DeleteByChainId(chain.Id)
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainDownload(ctx *gin.Context) {
	chainId, err := strconv.Atoi(ctx.Query("chainId"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "get_chainId_error", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = chainId
	isSuccess, chain := a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	reader, contentLength, name := a.chainService.DownloadChainArtifacts(chain)
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, name),
	}

	ctx.DataFromReader(http.StatusOK, contentLength, "application/x-tar", reader, extraHeaders)

}

func (a *ApiController) ChainPodsQuery(ctx *gin.Context) {
	chainId, err := strconv.Atoi(ctx.Query("chainId"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "get_chainId_error", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = chainId
	isSuccess, chain := a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	isSuccess, dat := a.chainService.QueryChainPods(chain)
	if isSuccess {
		gintool.ResultOk(ctx, dat)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "query_chain_error", nil, common.Lang))
	}

}

func (a *ApiController) ChangeChainResouces(ctx *gin.Context) {
	resouces := new(model.Resources)
	if err := ctx.ShouldBindJSON(resouces); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, dat := a.chainService.ChangeChainResouces(resouces)
	if isSuccess {
		gintool.ResultOk(ctx, dat)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "change_re_fail", nil, common.Lang))
	}

}
