package common

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/syyongx/ii18n"
)

var Lang string
var rLock sync.RWMutex

func I18nConfig() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Lang = ctx.GetHeader("X-Language")
		if Lang=="" {
			Lang="zh"
		}
		rLock.Lock()
		config := map[string]ii18n.Config{
			"app": ii18n.Config{
				SourceNewFunc: ii18n.NewJSONSource,
				BasePath:      "./common/locale",
				FileMap: map[string]string{
					"app": "app.json",
				},
			},
		}
		defer rLock.Unlock()
		ii18n.NewI18N(config)
		ctx.Next()
	}
}
