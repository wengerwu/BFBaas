package service

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	tsgutils "github.com/typa01/go-utils"

	"github.com/go-xorm/xorm"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-core/common/json"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

type ChaincodeService struct {
	DbEngine      *xorm.Engine
	FabircService *FabricService
}

func (l *ChaincodeService) Add(cc *entity.Chaincode) (bool, string) {
	cc.Created = time.Now().Unix()
	cc.Version = "1"
	cc.Status = 0
	cc.Secret = tsgutils.GUID()
	i, err := l.DbEngine.Insert(cc)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, ii18n.T("app", "add_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "add_fail", nil, common.Lang)
}

func (l *ChaincodeService) Update(cc *entity.Chaincode) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", cc.Id).Update(cc)
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "update_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "update_fail", nil, common.Lang)
}

func (l *ChaincodeService) Delete(id int) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&entity.Chaincode{})
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "delete_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "delete_fail", nil, common.Lang)
}

func (l *ChaincodeService) GetByChaincode(cc *entity.Chaincode) (bool, *entity.Chaincode) {
	has, err := l.DbEngine.Get(cc)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, cc
}

func (l *ChaincodeService) GetList(cc *entity.Chaincode, page, size int) (bool, []*entity.Chaincode, int64) {
	pager := gintool.CreatePager(page, size)
	ccs := make([]*entity.Chaincode, 0)

	values := make([]interface{}, 0)
	where := "1=1"
	if cc.ChaincodeName != "" {
		where += " and chaincode_name = ? "
		values = append(values, cc.ChaincodeName)
	}
	if cc.ChannelId != 0 {
		where += " and channel_id = ? "
		values = append(values, cc.ChannelId)
	}

	err := l.DbEngine.Where(where, values...).Limit(pager.PageSize, pager.NumStart).Find(&ccs)
	if err != nil {
		logger.Error(err.Error())
	}
	total, err := l.DbEngine.Where(where, values...).Count(new(entity.Chaincode))
	if err != nil {
		logger.Error(err.Error())
	}
	return true, ccs, total
}

func (l *ChaincodeService) GetAllList(chainId int) (bool, []*entity.Chaincode) {
	ccs := make([]*entity.Chaincode, 0)
	err := l.DbEngine.Where("channel_id = ?", chainId).Find(&ccs)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Println("ccs", ccs)
	return true, ccs
}

//判断是否存在相同的链码名称
func (l *ChaincodeService) HasChainCodeName(chain_id int, ccName string) bool {
	values := make([]interface{}, 0)
	values = append(values, chain_id)
	values = append(values, ccName)
	count, err := l.DbEngine.Where("chain_id = ? and chaincode_name = ?", values...).Count(new(entity.Chaincode))
	if count > 0 || err != nil {
		return true
	}
	return false
}

func (l *ChaincodeService) AddChaincode(chain *entity.Chain, channel *entity.Channel, cc *entity.Chaincode) (bool, string) {
	bys, err := ioutil.ReadFile(cc.GithubPath)
	if err != nil {
		logger.Error(err.Error())
		return false, ii18n.T("app", "add_fail", nil, common.Lang)
	}
	cc.ChainId = channel.ChainId
	cc.Version = "1"
	fc := entity.ParseFabircChannel(entity.ParseFabircChainAndChannel(chain, channel), cc)
	fc.ChaincodeBytes = bys
	resp := l.FabircService.UploadChaincode(fc)
	var ret gintool.RespData
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		fmt.Println(err)
		return false, ii18n.T("app", "add_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		cc.GithubPath = ret.Data.(string)
		cc.Created = time.Now().Unix()
		cc.Status = 0
		return l.Add(cc)
	} else {
		return false, ii18n.T("app", "add_fail", nil, common.Lang)
	}
}

func (l *ChaincodeService) DeployChaincode(chain *entity.Chain, channel *entity.Channel, cc *entity.Chaincode) (bool, string) {
	fc := entity.ParseFabircChannel(entity.ParseFabircChainAndChannel(chain, channel), cc)
	args := make([][]byte, 1)
	args[0] = []byte("init")

	for _, v := range strings.Split(cc.Args, ",") {
		args = append(args, []byte(v))
	}
	fc.Args = args
	resp := l.FabircService.BuildChaincode(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "deploy_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		cc.Status = 1
		return l.Update(cc)
	} else {
		return false, ii18n.T("app", "deploy_fail", nil, common.Lang)
	}
}

func (l *ChaincodeService) UpgradeChaincode(chain *entity.Chain, channel *entity.Channel, cc *entity.Chaincode) (bool, string) {
	bys, err := ioutil.ReadFile(cc.GithubPath)
	if err != nil {
		return false, ii18n.T("app", "upgrade_fail", nil, common.Lang)
	}
	v, err := strconv.Atoi(cc.Version)
	if err != nil {
		return false, ii18n.T("app", "version_error", nil, common.Lang)
	}
	cc.Version = fmt.Sprintf("%d", v+1)

	fc := entity.ParseFabircChannel(entity.ParseFabircChainAndChannel(chain, channel), cc)
	fc.ChaincodeBytes = bys
	resp := l.FabircService.UploadChaincode(fc)
	var ret gintool.RespData
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "upgrade_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		cc.GithubPath = ret.Data.(string)
		fc.ChaincodePath = ret.Data.(string)
	} else {
		return false, ii18n.T("app", "upload_fail", nil, common.Lang)
	}

	args := make([][]byte, 1)
	args[0] = []byte("init")
	for _, v := range strings.Split(cc.Args, ",") {
		args = append(args, []byte(v))
	}
	fc.Args = args
	resp = l.FabircService.UpdateChaincode(fc)
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "upgrade_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		return l.Update(cc)
	} else {
		return false, ii18n.T("app", "upgrade_fail", nil, common.Lang)
	}
}

func (l *ChaincodeService) InvokeChaincode(chain *entity.Chain, channel *entity.Channel, cc *entity.Chaincode) (bool, string) {
	fc := entity.ParseFabircChannel(entity.ParseFabircChainAndChannel(chain, channel), cc)
	args := make([][]byte, 0)

	for _, v := range strings.Split(cc.Args, ",") {
		args = append(args, []byte(v))
	}
	fc.Args = args
	resp := l.FabircService.InvokeChaincode(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		fmt.Println("resp", string(resp))
		fmt.Println(ret)
		fmt.Println(err)
		return false, ii18n.T("app", "invoke_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		return true, ret.Data.(string)
	} else {
		return false, ii18n.T("app", "invoke_fail", nil, common.Lang)
	}
}

func (l *ChaincodeService) QueryChaincode(chain *entity.Chain, channel *entity.Channel, cc *entity.Chaincode) (bool, string) {
	fc := entity.ParseFabircChannel(entity.ParseFabircChainAndChannel(chain, channel), cc)
	args := make([][]byte, 0)

	for _, v := range strings.Split(cc.Args, ",") {
		args = append(args, []byte(v))
	}
	fc.Args = args
	resp := l.FabircService.QueryChaincode(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "query_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		return true, ret.Data.(string)
	} else {
		return false, ii18n.T("app", "query_fail", nil, common.Lang)
	}
}

func (l *ChaincodeService) QueryLedger(chain *entity.Chain, channel *entity.Channel) (bool, interface{}) {
	fc := entity.ParseFabircChainAndChannel(chain, channel)
	resp := l.FabircService.QueryLedger(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "query_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		return true, ret.Data
	} else {
		return false, ret.Msg
	}

}

func (l *ChaincodeService) QueryLatestBlocks(chain *entity.Chain, channel *entity.Channel) (bool, interface{}) {
	fc := entity.ParseFabircChainAndChannel(chain, channel)
	resp := l.FabircService.QueryLatestBlocks(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "query_fail", nil, common.Lang)
	}
	if ret.Code == 0 {
		return true, ret.Data
	} else {
		return false, ret.Msg
	}
}

func (l *ChaincodeService) QueryBlock(chain *entity.Chain, channel *entity.Channel, search string) (bool, interface{}) {
	fc := entity.ParseFabircChainAndChannel(chain, channel)
	resp := l.FabircService.QueryBlock(fc, search)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "query_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		return true, ret.Data
	} else {
		return false, ret.Msg
	}

}

func NewChaincodeService(engine *xorm.Engine, fabircService *FabricService) *ChaincodeService {
	return &ChaincodeService{
		DbEngine:      engine,
		FabircService: fabircService,
	}
}
