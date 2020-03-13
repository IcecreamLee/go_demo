package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"os"
)

var sessionStore *sessions.CookieStore

func init() {
	sessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

func GetSession(c *gin.Context) *sessions.Session {
	session, _ := sessionStore.Get(c.Request, "SID")
	return session
}
