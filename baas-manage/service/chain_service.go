package service

import (
	"bytes"
	"io"

	"github.com/go-xorm/xorm"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-core/common/json"
	"github.com/paybf/baasmanager/baas-core/core/model"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

type ChainService struct {
	DbEngine      *xorm.Engine
	FabircService *FabricService
	UserService   *UserService
}

//判断同一个userid链下边有没有相同的链名
func (l *ChainService) IsHaveName(chain *entity.Chain, isAdd bool) bool {
	values := make([]interface{}, 0)
	where := "open_id = ? and name = ?"
	values = append(values, chain.OpenId)
	values = append(values, chain.Name)
	if !isAdd {
		where += " and id != ?"
		values = append(values, chain.Id)
	}
	count, err := l.DbEngine.Where(where, values...).Count(new(entity.Chain))
	if count > 0 || err != nil {
		return true
	}
	return false
}

func (l *ChainService) Add(chain *entity.Chain) (bool, string) {
	i, err := l.DbEngine.Insert(chain)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, ii18n.T("app", "add_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "add_fail", nil, common.Lang)
}

func (l *ChainService) Update(chain *entity.Chain) (bool, string) {

	i, err := l.DbEngine.Where("id = ?", chain.Id).Update(chain)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, ii18n.T("app", "update_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "update_fail", nil, common.Lang)
}

func (l *ChainService) UpdateStatus(chain *entity.Chain) (bool, string) {
	sql := "update chain set status = ? where id = ?"
	res, err := l.DbEngine.Exec(sql, chain.Status, chain.Id)
	if err != nil {
		logger.Error(err.Error())
	}

	r, err := res.RowsAffected()
	if err == nil && r > 0 {
		return true, ii18n.T("app", "update_success", nil, common.Lang)
	}

	return false, ii18n.T("app", "update_fail", nil, common.Lang)
}

func (l *ChainService) Delete(id int) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&entity.Chain{})
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "delete_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "delete_fail", nil, common.Lang)
}

func (l *ChainService) GetByChain(chain *entity.Chain) (bool, *entity.Chain) {
	has, err := l.DbEngine.Get(chain)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, chain
}

func (l *ChainService) GetList(chain *entity.Chain, page, size int) (bool, []*entity.Chain, int64) {
	pager := gintool.CreatePager(page, size)
	chains := make([]*entity.Chain, 0)

	values := make([]interface{}, 0)
	where := "1=1"
	if chain.Name != "" {
		where += " and name = ? "
		values = append(values, chain.Name)
	}
	if chain.OpenId != "" {
		where += " and open_id = ? "
		values = append(values, chain.OpenId)
	}
	if chain.Consensus != "" {
		where += " and consensus = ? "
		values = append(values, chain.Consensus)
	}
	if chain.PeersOrgs != "" {
		where += " and peers_orgs like ? "
		values = append(values, "%"+chain.PeersOrgs+"%")
	}
	if chain.TlsEnabled != "" {
		where += " and tls_enabled = ? "
		values = append(values, chain.TlsEnabled)
	}

	err := l.DbEngine.Where(where, values...).Limit(pager.PageSize, pager.NumStart).Find(&chains)
	if err != nil {
		logger.Error(err.Error())
	}

	for i := 0; i < len(chains); i++ {
		user := &entity.User{OpenId: chains[i].OpenId}
		_, user = l.UserService.GetByUser(user)
		chains[i].UserName = user.Name
		chains[i].UserPhone = user.Phone
	}

	total, err := l.DbEngine.Where(where, values...).Count(new(entity.Chain))
	if err != nil {
		logger.Error(err.Error())
	}
	return true, chains, total
}

func (l *ChainService) GetUserIdList(userNameOrPhone string) []entity.UserIdAccount {
	users := make([]entity.User, 0)
	err := l.DbEngine.Where("user_type=1 and name LIKE '%" + userNameOrPhone + "%' or phone LIKE '" + userNameOrPhone + "%'").Find(&users)
	if err != nil {
		logger.Error(err.Error())
	}

	useridList := make([]entity.UserIdAccount, 0)
	for _, v := range users {
		useridList = append(useridList, entity.UserIdAccount{OpenId: v.OpenId, Account: v.Name + "/" + v.Phone})
	}
	return useridList
}

func (l *ChainService) BuildChain(chain *entity.Chain) (bool, string) {
	fc := entity.ParseFabircChain(chain)
	resp := l.FabircService.DefChain(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "build_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		chain.Status = 1
		return l.UpdateStatus(chain)
	} else {
		return false, ret.Msg
	}

}

func (l *ChainService) RunChain(chain *entity.Chain) (bool, string) {
	fc := entity.ParseFabircChain(chain)
	resp := l.FabircService.DeployK8sData(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "run_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		chain.Status = 2
		return l.UpdateStatus(chain)
	} else {
		return false, ret.Msg
	}

}

func (l *ChainService) QueryChainPods(chain *entity.Chain) (bool, interface{}) {
	fc := entity.ParseFabircChain(chain)
	resp := l.FabircService.QueryChainPods(fc)
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

func (l *ChainService) ChangeChainResouces(resouce *model.Resources) (bool, interface{}) {
	resp := l.FabircService.ChangeChainPodResources(*resouce)
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

func (l *ChainService) StopChain(chain *entity.Chain) (bool, string) {
	fc := entity.ParseFabircChain(chain)
	resp := l.FabircService.StopChain(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "run_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		chain.Status = 3
		return l.UpdateStatus(chain)
	} else {
		return false, ii18n.T("app", "build_fail", nil, common.Lang)
	}

}

func (l *ChainService) ReleaseChain(chain *entity.Chain) (bool, string) {
	fc := entity.ParseFabircChain(chain)
	resp := l.FabircService.ReleaseChain(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, ii18n.T("app", "run_fail", nil, common.Lang)
	}

	if ret.Code == 0 {
		chain.Status = 0
		return l.UpdateStatus(chain)
	} else {
		return false, ii18n.T("app", "build_fail", nil, common.Lang)
	}

}

func (l *ChainService) DownloadChainArtifacts(chain *entity.Chain) (io.Reader, int64, string) {
	fc := entity.ParseFabircChain(chain)
	bts := l.FabircService.DownloadChainArtifacts(fc)
	reader := bytes.NewReader(bts)
	contentLength := reader.Len()
	return reader, int64(contentLength), chain.Name + ".tar"

}

func NewChainService(engine *xorm.Engine, fabircService *FabricService, userService *UserService) *ChainService {
	return &ChainService{
		DbEngine:      engine,
		FabircService: fabircService,
		UserService:   userService,
	}
}
