package controller

import (
	"strconv"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

func (a *ApiController) ChaincodeAdd(ctx *gin.Context) {
	cc := new(entity.Chaincode)
	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	ccNameRule, _ := regexp.MatchString("^[A-Za-z0-9]+$", cc.ChaincodeName)
	if !ccNameRule || cc.ChaincodeName == "" || cc.Policy == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "args_fail", nil, common.Lang))
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "channel_not_exist", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	hasChainCodeName := a.chaincodeService.HasChainCodeName(channel.ChainId, cc.ChaincodeName)
	if hasChainCodeName {
		gintool.ResultFail(ctx, ii18n.T("app", "hasChanCodeName", nil, common.Lang))
		return
	}

	cc.OpenId = chain.OpenId
	isSuccess, msg := a.chaincodeService.AddChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeDeploy(ctx *gin.Context) {
	cc := new(entity.Chaincode)
	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "channel_not_exist", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	isSuccess, msg := a.chaincodeService.DeployChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeUpgrade(ctx *gin.Context) {
	cc := new(entity.Chaincode)
	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "channel_not_exist", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	isSuccess, msg := a.chaincodeService.UpgradeChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeQuery(ctx *gin.Context) {
	cc := new(entity.Chaincode)
	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "channel_not_exist", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	isSuccess, msg := a.chaincodeService.QueryChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeLedgerQuery(ctx *gin.Context) {
	channelId, err := strconv.Atoi(ctx.Query("channelId"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "channelId_error", nil, common.Lang))
		return
	}

	channel := new(entity.Channel)
	channel.Id = channelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "channel_not_exist", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	isSuccess, msg := a.chaincodeService.QueryLedger(chain, channel)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeLatestBlocksQuery(ctx *gin.Context) {
	channelId, err := strconv.Atoi(ctx.Query("channelId"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "channelId_error", nil, common.Lang))
		return
	}

	channel := new(entity.Channel)
	channel.Id = channelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "channel_not_exist", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	isSuccess, msg := a.chaincodeService.QueryLatestBlocks(chain, channel)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeBlockQuery(ctx *gin.Context) {
	channelId, err := strconv.Atoi(ctx.Query("channelId"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "channelId_error", nil, common.Lang))
		return
	}
	search := ctx.Query("search")

	channel := new(entity.Channel)
	channel.Id = channelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "channel_not_exist", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	isSuccess, msg := a.chaincodeService.QueryBlock(chain, channel, search)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeInvoke(ctx *gin.Context) {
	cc := new(entity.Chaincode)
	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "channel_not_exist", nil, common.Lang))
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, ii18n.T("app", "chain_not_exist", nil, common.Lang))
		return
	}

	isSuccess, msg := a.chaincodeService.InvokeChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeGet(ctx *gin.Context) {
	cc := new(entity.Chaincode)
	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, chain := a.chaincodeService.GetByChaincode(cc)
	if isSuccess {
		gintool.ResultOk(ctx, chain)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "query_chaincode_error", nil, common.Lang))
	}
}

func (a *ApiController) ChaincodeList(ctx *gin.Context) {
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

	name := ctx.Query("chaincodeName")
	channelId, err := strconv.Atoi(ctx.Query("channelId"))
	if err != nil {
		gintool.ResultFail(ctx, ii18n.T("app", "channelId_error", nil, common.Lang))
		return
	}

	b, list, total := a.chaincodeService.GetList(&entity.Chaincode{
		ChaincodeName: name,
		ChannelId:     channelId,
	}, page, limit)

	if b {
		gintool.ResultList(ctx, list, total)

	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "query_chaincode_list_error", nil, common.Lang))
	}
}

func (a *ApiController) ChaincodeUpdate(ctx *gin.Context) {
	cc := new(entity.Chaincode)
	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chaincodeService.Update(cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeDeleted(ctx *gin.Context) {
	cc := new(entity.Chaincode)
	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chaincodeService.Delete(cc.Id)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}
