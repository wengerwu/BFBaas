package service

import (
	"github.com/go-xorm/xorm"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

type BaasComboService struct {
	DbEngine *xorm.Engine
}

func (l *BaasComboService) Add(entity *entity.BaasCombo) (bool, string) {
	i, err := l.DbEngine.Insert(entity)
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "add_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "add_fail", nil, common.Lang)
}

func (l *BaasComboService) Update(entity *entity.BaasCombo) (bool, string) {
	res, err := l.DbEngine.Exec("update baas_combo set displayName = ?, sort = ?, price = ?, remark = ? where id = ?", entity.DisplayName, entity.Sort, entity.Price, entity.Remark, entity.Id)
	row, err := res.RowsAffected()
	if err != nil {
		logger.Error(err.Error())
		return false, ii18n.T("app", "update_fail", nil, common.Lang)
	}

	if row > 0 {
		return true, ii18n.T("app", "update_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "update_fail", nil, common.Lang)
}

/*删除指定的套餐*/
func (l *BaasComboService) Delete(id string) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&entity.BaasCombo{})
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "delete_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "delete_fail", nil, common.Lang)
}

/*
获取所有套餐，无需分页
*/
func (l *BaasComboService) GetAllList() (bool, []entity.BaasCombo) {
	baasCombos := make([]entity.BaasCombo, 0)
	err := l.DbEngine.Asc("sort").Find(&baasCombos)
	if err != nil {
		logger.Error(err.Error())
		return false, nil
	}
	return true, baasCombos
}

func (l *BaasComboService) GetCombo(displayName string) (bool, []entity.BaasCombo) {
	baasCombos := make([]entity.BaasCombo, 0)
	err := l.DbEngine.Where("displayName = ?", displayName).Find(&baasCombos)
	if err != nil {
		logger.Error(err.Error())
		return false, nil
	}
	return true, baasCombos
}

func NewComboService(engine *xorm.Engine) *BaasComboService {
	return &BaasComboService{
		DbEngine: engine,
	}
}
