package service

import (
	"fmt"
	"time"

	"github.com/go-xorm/core"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-xorm/xorm"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	jwttool "github.com/paybf/baasmanager/baas-core/common/jwt"
	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
	tsgutils "github.com/typa01/go-utils"
)

const TokenKey = "baas user secret"

type UserService struct {
	DbEngine *xorm.Engine
}

func (l *UserService) isAccountHave(user *entity.User, isAdd bool) bool {
	values := make([]interface{}, 0)
	where := "account = ?"
	values = append(values, user.Account)
	if !isAdd {
		where += " and id != ?"
		values = append(values, user.Id)
	}
	i, err := l.DbEngine.Where(where, values...).Count(new(entity.User))
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true
	}
	return false
}

func (l *UserService) Add(user *entity.User) (bool, string) {
	ishave := l.isAccountHave(user, true)
	if ishave {
		return false, ii18n.T("app", "accountHave", nil, common.Lang)
	}
	user.OpenId = tsgutils.GUID()
	i, err := l.DbEngine.Insert(user)
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "add_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "add_fail", nil, common.Lang)
}

func (l *UserService) Update(user *entity.User) (bool, string) {
	ishave := l.isAccountHave(user, false)
	if ishave {
		return false, ii18n.T("app", "accountHave", nil, common.Lang)
	}

	i, err := l.DbEngine.Where("id = ?", user.Id).Update(user)
	if err != nil {
		logger.Error(err.Error())
	}

	if i == 1 || i == 0 {
		return true, ii18n.T("app", "update_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "update_fail", nil, common.Lang)
}

func (l *UserService) UpdateByAccount(user *entity.User) (bool, string) {
	i, err := l.DbEngine.Where("account = ?", user.Account).Update(user)
	if err != nil {
		logger.Error(err.Error())
	}

	if i == 1 || i == 0 {
		return true, ii18n.T("app", "update_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "update_fail", nil, common.Lang)
}

func (l *UserService) UpdateByPhone(user *entity.User) (bool, string) {
	i, err := l.DbEngine.Where("phone = ?", user.Phone).Update(user)
	if err != nil {
		logger.Error(err.Error())
	}

	if i == 1 || i == 0 {
		return true, ii18n.T("app", "update_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "update_fail", nil, common.Lang)
}

func (l *UserService) Delete(Userid string) (bool, string) {
	ur := new(entity.UserRole)
	ur.UserId = Userid
	l.DelAuth(ur)

	i, err := l.DbEngine.Where("userid = ?", Userid).Delete(&entity.User{})
	if err != nil {
		logger.Error(err.Error())
	}

	if i > 0 {
		return true, ii18n.T("app", "delete_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "delete_fail", nil, common.Lang)
}

func (l *UserService) GetByUser(user *entity.User) (bool, *entity.User) {
	has, err := l.DbEngine.Get(user)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, user
}

func (l *UserService) GetUserByPhone(phone string) (bool, *entity.User) {
	user := new(entity.User)
	has, err := l.DbEngine.Where("phone = ?", phone).Get(user)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, user
}

func (l *UserService) GetUserByAccount(account string) (bool, *entity.User) {
	user := new(entity.User)
	has, err := l.DbEngine.Where("account = ?", account).Get(user)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, user
}

func (l *UserService) GetList(user *entity.User, page, size int) (bool, []entity.UserDetail, int64) {
	pager := gintool.CreatePager(page, size)
	users := make([]*entity.User, 0)
	values := make([]interface{}, 0)
	where := "1=1"
	if user.Account != "" {
		where += " and account = ? "
		values = append(values, user.Account)
	}
	if user.Name != "" {
		where += " and name like ? "
		values = append(values, "%"+user.Name+"%")
	}
	where += " and user_type=0"

	err := l.DbEngine.Where(where, values...).Limit(pager.PageSize, pager.NumStart).Find(&users)
	if err != nil {
		logger.Error(err.Error())
	}

	total, err := l.DbEngine.Where(where, values...).Count(new(entity.User))
	if err != nil {
		logger.Error(err.Error())
	}

	userIds := make([]string, len(users))
	userDatas := make([]entity.UserDetail, len(users))
	for i, u := range users {
		userIds[i] = u.Userid
		userDatas[i].Id = u.Id
		userDatas[i].Userid = u.Userid
		userDatas[i].Account = u.Account
		userDatas[i].Password = u.Password
		userDatas[i].Avatar = u.Avatar
		userDatas[i].Name = u.Name
		userDatas[i].Created = u.Created
	}

	roles := make([]entity.UserRole, 0)
	err = l.DbEngine.In("user_id", userIds).Find(&roles)
	if err != nil {
		logger.Error(err.Error())
	}

	for i, d := range userDatas {
		keys := make([]string, 0)
		for _, r := range roles {
			if r.UserId == d.Userid {
				keys = append(keys, r.RoleKey)
			}
		}
		d.Roles = keys
		userDatas[i] = d
	}

	return true, userDatas, total
}

func (l *UserService) GetToken(user *entity.User) *entity.JwtToken {
	info := make(map[string]interface{})
	now := time.Now()
	info["userId"] = user.Id
	info["exp"] = now.Add(time.Hour * 1).Unix() // 1 小时过期
	info["iat"] = now.Unix()
	tokenString := jwttool.CreateToken(TokenKey, info)
	return &entity.JwtToken{
		Token: tokenString,
	}
}

func (l *UserService) CheckToken(token string, user *entity.User) (*entity.UserInfo, error) {
	info, ok := jwttool.ParseToken(token, TokenKey)
	infoMap := info.(jwt.MapClaims)
	if ok {
		expTime := infoMap["exp"].(float64)
		if float64(time.Now().Unix()) >= expTime {
			return nil, fmt.Errorf("%s", ii18n.T("app", "Token_timeout", nil, common.Lang))
		} else {
			l.DbEngine.Get(user)
			ur := make([]entity.UserRole, 0)
			err := l.DbEngine.Where("user_id = ?", user.Userid).Find(&ur)
			if err != nil {
				logger.Error(err.Error())
			}
			roles := make([]string, len(ur))
			for i, m := range ur {
				roles[i] = m.RoleKey
			}
			info := &entity.UserInfo{
				Avatar:  user.Avatar,
				Roles:   roles,
				Name:    user.Name,
				Account: user.Account,
			}
			return info, nil
		}
	} else {
		return nil, fmt.Errorf("%s", ii18n.T("app", "token_not_exist", nil, common.Lang))
	}
}

func (l *UserService) AddAuth(ur *entity.UserRole) (bool, string) {
	i, err := l.DbEngine.ID(core.PK{ur.UserId}).Insert(ur)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, ii18n.T("app", "add_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "add_fail", nil, common.Lang)
}

func (l *UserService) UpdateAuth(ur *entity.UserRole) (bool, string) {
	values := make([]interface{}, 0)
	where := "user_id = ?"
	values = append(values, ur.UserId)
	i, err := l.DbEngine.Where(where, values...).Update(ur)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, ii18n.T("app", "update_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "update_fail", nil, common.Lang)
}

func (l *UserService) DelAuth(ur *entity.UserRole) (bool, string) {
	values := make([]interface{}, 0)
	where := "user_id = ?"
	values = append(values, ur.UserId)
	i, err := l.DbEngine.Where(where, values...).Delete(&entity.UserRole{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, ii18n.T("app", "delete_success", nil, common.Lang)
	}
	return false, ii18n.T("app", "delete_fail", nil, common.Lang)
}

func NewUserService(engine *xorm.Engine) *UserService {
	return &UserService{
		DbEngine: engine,
	}
}
