package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paybf/baasmanager/baas-core/common/gintool"
	"github.com/paybf/baasmanager/baas-gateway/config"
	"github.com/paybf/baasmanager/baas-gateway/entity"
	"github.com/syyongx/ii18n"
)

func CheckSign(timestamp int64, sign string, src string) gintool.RespData {
	currTimestamp := time.Now().UnixNano() / 1e6
	if currTimestamp-timestamp > 120*1000 {
		return gintool.RespData{Code: gintool.Fail, Data: ii18n.T("app", "sign_expire", nil, Lang)}
	}

	currSign := GetSign(src)
	if sign != currSign {
		return gintool.RespData{Code: gintool.Fail, Data: ii18n.T("app", "sign_invalid", nil, Lang)}
	}

	return gintool.RespData{Code: gintool.Success}
}

func GetSign(src string) string {
	hashSuffix := config.Config.GetString("IdentityServer.HashSuffix")
	s := md5.New()
	s.Write([]byte(src + "&" + hashSuffix))
	return hex.EncodeToString(s.Sum(nil))
}

func GetSubjectIdValue(ctx *gin.Context) (string, error) {
	cookie, _ := ctx.Cookie("Sys-SubjectId")
	return cookie, nil
}

func GetUserTokenInfoByGuid(guid string) (*entity.UserTokenInfo, error) {
	val, err := RedisClient.Get(guid).Result()
	if err != nil {
		return nil, err
	}
	userTokenInfo := new(entity.UserTokenInfo)
	json.Unmarshal([]byte(val), userTokenInfo)
	return userTokenInfo, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
