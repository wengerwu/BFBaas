package controller

import (
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

func (a *ApiController) ChannelAdd(ctx *gin.Context) {
	channel := new(entity.Channel)
	if err := ctx.ShouldBindJSON(channel); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channelNameRule, _ := regexp.MatchString("^[a-z][a-zA-Z0-9]+$", channel.ChannelName)
	if channel.ChannelName == "" || !channelNameRule || channel.Orgs == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "args_fail", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain := a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	hasChannelName := a.channelService.HasChannelName(chain.Id, channel.ChannelName)
	if hasChannelName {
		gintool.ResultFail(ctx, ii18n.T("app", "hasChannelName", nil, common.Lang))
		return
	}

	channel.OpenId = chain.OpenId
	isSuccess, msg := a.channelService.AddChannel(chain, channel)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChannelGet(ctx *gin.Context) {
	chn := new(entity.Channel)
	if err := ctx.ShouldBindJSON(chn); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, chn := a.channelService.GetByChannel(chn)
	if isSuccess {
		gintool.ResultOk(ctx, chn)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "channel_not_exist", nil, common.Lang))
	}
}

func (a *ApiController) ChannelAll(ctx *gin.Context) {
	chainId, err := strconv.Atoi(ctx.Query("chainId"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "get_chainId_error", nil, common.Lang))
		return
	}
	isSuccess, data := a.channelService.GetAllList(chainId)
	if isSuccess {
		gintool.ResultOk(ctx, data)
	} else {
		gintool.ResultFail(ctx, data)
	}
}
