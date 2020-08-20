package service

import (
	"encoding/json"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

type ChannelService struct {
	DbEngine      *xorm.Engine
	FabircService *FabricService
}

func (l *ChannelService) Add(channel *entity.Channel) (bool, string) {
	i, err := l.DbEngine.Insert(channel)
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "add_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "add_fail", nil, common.Lang)
}

func (l *ChannelService) Update(channel *entity.Channel) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", channel.Id).Update(channel)
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "update_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "update_fail", nil, common.Lang)
}

func (l *ChannelService) Delete(id int) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&entity.Channel{})
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "delete_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "delete_fail", nil, common.Lang)
}

func (l *ChannelService) DeleteByChainId(id int) (bool, string) {
	sql := "delete from chaincode where channel_id in ( select id from channel where chain_id = ?)"
	_, err := l.DbEngine.Exec(sql, id)
	if err != nil {
		logger.Error(err.Error())
	}

	i, err := l.DbEngine.Where("chain_id = ?", id).Delete(&entity.Channel{})
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "delete_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "delete_fail", nil, common.Lang)
}

func (l *ChannelService) GetByChannel(channel *entity.Channel) (bool, *entity.Channel) {
	has, err := l.DbEngine.Get(channel)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, channel
}

func (l *ChannelService) GetList(channel *entity.Channel, page, size int) (bool, []*entity.Channel) {
	channels := make([]*entity.Channel, 0)
	values := make([]interface{}, 0)
	where := "1=1"

	err := l.DbEngine.Where(where, values...).Limit(size, page).Find(&channels)
	if err != nil {
		logger.Error(err.Error())
	}

	return true, channels
}

func (l *ChannelService) GetAllList(chainId int) (bool, []*entity.Channel) {
	channels := make([]*entity.Channel, 0)
	err := l.DbEngine.Where("chain_id = ?", chainId).Find(&channels)
	if err != nil {
		logger.Error(err.Error())
	}

	return true, channels
}

func (l *ChannelService) AddChannel(chain *entity.Chain, channel *entity.Channel) (bool, string) {
	fc := entity.ParseFabircChainAndChannel(chain, channel)
	resp := l.FabircService.DefChannel(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "add_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		channel.Created = time.Now().Unix()
		return l.Add(channel)
	} else {
		return false, ii18n.T("app", "add_fail", nil, common.Lang)
	}
}

//判断是否同一个链下边有相同的名称
func (l *ChannelService) HasChannelName(channelId int, channelName string) bool {
	values := make([]interface{}, 0)
	values = append(values, channelId)
	values = append(values, channelName)
	count, err := l.DbEngine.Where("chain_id = ? and channel_name = ?", values...).Count(new(entity.Channel))
	if count > 0 || err != nil {
		return true
	}
	return false
}

func NewChannelService(engine *xorm.Engine, fabircService *FabricService) *ChannelService {
	return &ChannelService{
		DbEngine:      engine,
		FabircService: fabircService,
	}
}
