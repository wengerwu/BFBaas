package service

import (
	"github.com/go-xorm/xorm"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

type BaasOrderService struct {
	DbEngine *xorm.Engine
}

func (l *BaasOrderService) Add(entity *entity.BaasOrder) (bool, string) {
	i, err := l.DbEngine.Insert(&entity)
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "add_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "add_fail", nil, common.Lang)
}

func (l *BaasOrderService) Delete(id string) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&entity.BaasOrder{})
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "delete_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "delete_fail", nil, common.Lang)
}

func (l *BaasOrderService) Update(entity *entity.BaasOrder) (bool, string) {
	if entity.Price == 0 {
		return false, ii18n.T("app", "Data_cannot_be_Zero", nil, common.Lang)
	}
	if entity.STATUS == 0 {
		entity.STATUS = 1
	}
	if entity.STATUS == 2 {
		entity.STATUS = 3
	}
	i, err := l.DbEngine.Where("id = ?", entity.Id).Update(entity)
	if err != nil {
		logger.Error(err.Error())
		return false, err.Error()
	}

	if i > 0 {
		return true, ii18n.T("app", "update_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "update_fail", nil, common.Lang)
}

func (l *BaasOrderService) GetList(status string, baasOrder *entity.BaasOrder, page, size int) (bool, []*entity.BaasOrder, int64) {
	pager := gintool.CreatePager(page, size)
	baasOrders := make([]*entity.BaasOrder, 0)
	values := make([]interface{}, 0)
	where := "1=1"
	if status != "" {
		where += " and STATUS = ? "
		values = append(values, status)
	}
	if baasOrder.TagName != "" {
		where += " and TagName = ? "
		values = append(values, baasOrder.TagName)
	}
	where += " and IsDelete <> 1 "
	err := l.DbEngine.Where(where, values...).Limit(pager.PageSize, pager.NumStart).Find(&baasOrders)
	if err != nil {
		logger.Error(err.Error())
		return false, nil, 0
	}

	total, err := l.DbEngine.Where(where, values...).Count(new(entity.BaasOrder))
	if err != nil {
		logger.Error(err.Error())
		return false, nil, 0
	}
	return true, baasOrders, total
}

func NewOrderService(engine *xorm.Engine) *BaasOrderService {
	return &BaasOrderService{
		DbEngine: engine,
	}
}
