package controller

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"github.com/paybf/baasmanager/baas-gateway/common"
	"github.com/paybf/baasmanager/baas-gateway/config"
	"github.com/syyongx/ii18n"

	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	tsgutils "github.com/typa01/go-utils"
)

func (a *ApiController) UserAdd(ctx *gin.Context) {
	user := new(entity.User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if user.Account == "" || user.Name == "" || user.Password == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "notValue", nil, common.Lang))
		return
	}

	hashSuffix := config.Config.GetString("IdentityServer.HashSuffix")

	cryto := md5.New()
	cryto.Write([]byte(user.Password + hashSuffix))
	pwd := hex.EncodeToString(cryto.Sum(nil))

	user.Password = pwd
	user.Created = time.Now().Unix()
	user.Userid = tsgutils.GUID()
	user.UserType = 0
	isSuccess, msg := a.userService.Add(user)
	if !isSuccess {
		gintool.ResultFail(ctx, msg)
		return
	}

	ur := new(entity.UserRole)
	ur.UserId = user.Userid
	ur.RoleKey = "user"
	isSuccess, msg = a.userService.AddAuth(ur)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) UserAddAuth(ctx *gin.Context) {
	ur := new(entity.UserRole)
	if err := ctx.ShouldBindJSON(ur); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.userService.AddAuth(ur)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) UserUpdateAuth(ctx *gin.Context) {
	ur := new(entity.UserRole)
	if err := ctx.ShouldBindJSON(ur); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if ur.RoleKey == "" || ur.UserId == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "notValue", nil, common.Lang))
		return
	}

	isSuccess, msg := a.userService.UpdateAuth(ur)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) UserDelAuth(ctx *gin.Context) {
	ur := new(entity.UserRole)
	if err := ctx.ShouldBindJSON(ur); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if ur.UserId == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "notValue", nil, common.Lang))
		return
	}

	isSuccess, msg := a.userService.DelAuth(ur)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) UserUpdate(ctx *gin.Context) {
	user := new(entity.User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if user.Account == "" || user.Name == "" || user.Password == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "notValue", nil, common.Lang))
		return
	}

	isSuccess, msg := a.userService.Update(user)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) UserDelete(ctx *gin.Context) {
	user := new(entity.User)
	if err := ctx.ShouldBindJSON(user); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if user.Userid == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "notValue", nil, common.Lang))
		return
	}

	isSuccess, msg := a.userService.Delete(user.Userid)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) IsAdmin(l *entity.LoginFrom) bool {
	u := new(entity.User)
	u.Account = l.Username
	has, user := a.userService.GetUserByAccount(u.Account)
	if has {
		if user.UserType == 0 {
			return false
		}
	}
	return true
}

func (a *ApiController) UserLogin(ctx *gin.Context) {
	login := new(entity.LoginFrom)
	if err := ctx.ShouldBind(&login); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if a.IsAdmin(login) {
		gintool.ResultFailData(ctx, ii18n.T("app", "You_have_no_authority", nil, common.Lang), nil)
		return
	}

	if login.Username == "" || login.Password == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "notValue", nil, common.Lang))
		return
	}

	params := "username=" + login.Username + "&password=" + login.Password + "&timestamp=" + strconv.FormatInt(login.Timestamp, 10)
	respData := common.CheckSign(login.Timestamp, login.Sign, params)
	if respData.Code != gintool.Success {
		gintool.ResultFailData(ctx, respData.Data, nil)
		return
	}

	user := &entity.User{
		Account: login.Username,
	}
	has, u := a.userService.GetUserByAccount(user.Account)
	if !has {
		gintool.ResultFail(ctx, ii18n.T("app", "username_error", nil, common.Lang))
		return
	}
	if login.Password != u.Password {
		gintool.ResultFail(ctx, ii18n.T("app", "password_error", nil, common.Lang))
		return
	}

	type UserInfo map[string]interface{}
	token := a.userService.GetToken(u)
	gintool.SetSession(ctx, token.Token, u.Account)
	_, err := ctx.Cookie("Sys-SubjectId")
	if err != nil {
		ctx.SetCookie(
			"Sys-SubjectId",
			token.Token,
			3600,
			"/",
			config.Config.GetString("Cookie.Domain"),
			false,
			false)
	}
	userTokenInfo := new(entity.UserTokenInfo)
	userTokenInfo.Token = token.Token
	userTokenInfo.Account = user.Account
	userTokenInfo.UserId = user.Userid
	infojson, _ := json.Marshal(userTokenInfo)
	var exp = 7200 * time.Second
	err = common.RedisClient.Set(user.Userid, string(infojson), exp).Err()
	if err != nil {
		gintool.ResultFailData(ctx, ii18n.T("app", "redis_set_fail", nil, common.Lang), err)
		return
	}
	gintool.ResultOk(ctx, token)
}

func (a *ApiController) UserLogout(ctx *gin.Context) {
	_, err := common.GetSubjectIdValue(ctx)
	if err != nil {
		gintool.ResultFailData(ctx, ii18n.T("app", "get_token_fail", nil, common.Lang), err)
		return
	}

	ctx.SetCookie(
		"Sys-SubjectId",
		"",
		-1,
		"/",
		config.Config.GetString("Cookie.Domain"),
		false,
		false,
	)
	gintool.ResultMsg(ctx, ii18n.T("app", "login_out", nil, common.Lang))
}

func (a *ApiController) UserAuthorize(ctx *gin.Context) {
	m := make(map[string]interface{})
	m["code"] = 2

	token, err := common.GetSubjectIdValue(ctx)
	if err != nil {
		m["msg"] = err.Error()
		gintool.ResultMap(ctx, m)
		ctx.Abort()
		return
	}

	session := gintool.GetSession(ctx, token)
	if nil == session {
		m["msg"] = ii18n.T("app", "token_not_exist", nil, common.Lang)
		gintool.ResultMap(ctx, m)
		ctx.Abort()
		return
	}
	_, err = a.userService.CheckToken(token, &entity.User{Userid: session.(string)})

	if err != nil {
		if err.Error() == ii18n.T("app", "Token_timeout", nil, common.Lang) || err.Error() == ii18n.T("app", "token_not_exist", nil, common.Lang) {
			m["msg"] = err.Error()
			gintool.ResultMap(ctx, m)
		} else {
			gintool.ResultFail(ctx, err.Error())
		}
		ctx.Abort()
		return
	} else {
		ctx.Next()
	}
}

func (a *ApiController) UserInfo(ctx *gin.Context) {
	m := make(map[string]interface{})
	m["code"] = 2
	token, err := common.GetSubjectIdValue(ctx)
	if err != nil {
		m["msg"] = err.Error()
		gintool.ResultMap(ctx, m)
		ctx.Abort()
		return
	}

	session := gintool.GetSession(ctx, token)
	if nil == session {
		gintool.ResultFail(ctx, ii18n.T("app", "token_not_exist", nil, common.Lang))
		return
	}

	_, u := a.userService.GetUserByAccount(session.(string))
	user, err := a.userService.CheckToken(token, &entity.User{Userid: u.Userid, Account: session.(string)})
	if err != nil {
		if err.Error() == ii18n.T("app", "Token_timeout", nil, common.Lang) || err.Error() == ii18n.T("app", "token_not_exist", nil, common.Lang) {
			m := make(map[string]interface{})
			m["code"] = 2
			m["msg"] = err.Error()
			gintool.ResultMap(ctx, m)
			return
		}
		gintool.ResultFail(ctx, err.Error())
	} else {
		gintool.ResultOk(ctx, user)
	}
}

func (a *ApiController) UserList(ctx *gin.Context) {
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
	b, list, total := a.userService.GetList(&entity.User{Name: name}, page, limit)

	if b {
		gintool.ResultList(ctx, list, total)
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "Get_User_Fail", nil, common.Lang))
	}
}

func (a *ApiController) ChangePwd(ctx *gin.Context) {
	pwdFrom := new(entity.ChangePwdFrom)
	if err := ctx.ShouldBindJSON(pwdFrom); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	if pwdFrom.OldPassword == "" || pwdFrom.OldPassword == "" || pwdFrom.NewPassword == "" {
		gintool.ResultFail(ctx, ii18n.T("app", "notValue", nil, common.Lang))
		return
	}

	guid, err := common.GetSubjectIdValue(ctx)
	account := gintool.GetSession(ctx, guid)
	if err != nil {
		return
	}

	has, user := a.userService.GetUserByAccount(account.(string))
	if !has {
		gintool.ResultFail(ctx, ii18n.T("app", "Get_User_Fail", nil, common.Lang))
		return
	}
	if pwdFrom.NewPassword == user.Password {
		gintool.ResultFail(ctx, ii18n.T("app", "pwdNotSame", nil, common.Lang))
		return
	}
	user.Password = pwdFrom.NewPassword
	isSuccess, msg := a.userService.Update(user)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ResetPwd(ctx *gin.Context) {
	u := new(entity.User)
	if err := ctx.ShouldBindJSON(u); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	has, user := a.userService.GetUserByAccount(u.Account)
	if !has {
		gintool.ResultFail(ctx, ii18n.T("app", "Get_User_Fail", nil, common.Lang))
		return
	}

	hashSuffix := config.Config.GetString("IdentityServer.HashSuffix")
	s := md5.New()
	s.Write([]byte("123456" + hashSuffix))
	cryptoPwd:=hex.EncodeToString(s.Sum(nil))

	user.Password = cryptoPwd
	isSuccess, _ := a.userService.Update(user)
	if isSuccess {
		gintool.ResultMsg(ctx, ii18n.T("app", "pwd_reset_success", nil, common.Lang))
		return
	} else {
		gintool.ResultFail(ctx, ii18n.T("app", "pwd_reset_fail", nil, common.Lang))
	}
}
