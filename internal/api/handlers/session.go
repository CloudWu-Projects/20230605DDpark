package api

import (
	"encoding/gob"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 用户会话信息结构
type UserSession struct {
	IsLogin    bool                   `json:"is_login"`
	OpenID     string                 `json:"openid"`
	UserInfo   map[string]interface{} `json:"user_info"`
	LastActive time.Time              `json:"last_active"`
}

var (
	MaxSessionDuration = 10 * time.Minute
)

func SetSession(r *gin.Engine) {

	gob.Register(UserSession{})
	store := cookie.NewStore([]byte("secret-key-123456")) // 生产环境使用更复杂的密钥
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   int(MaxSessionDuration.Seconds()), //86400 * 7, // 7天
		HttpOnly: true,
		Secure:   false, // 本地开发设为false，生产环境设为true
	})
	//store.RegisterType(UserSession{})
	r.Use(sessions.Sessions("qq_session", store))

}
