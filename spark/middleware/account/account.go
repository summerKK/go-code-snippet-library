package account

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/summerKK/go-code-snippet-library/session"
	"net/http"
)

func processRequest(c *gin.Context) {
	var redisSession session.ISession
	defer func() {
		if redisSession == nil {
			redisSession, _ = session.Session.Create()
		}
		c.Set(sparkSessionName, redisSession)
	}()
	c.Set(sparkUserId, int64(0))
	c.Set(sparkLoginStatus, 0)
	// 获取cookie
	sessionId, err := c.Cookie(cookieSessionId)
	if err != nil {
		return
	}
	redisSession, err = session.Session.Get(sessionId)
	if err != nil {
		return
	}
	tempUserId, err := redisSession.Get(sparkUserId)
	if err != nil {
		return
	}
	userId, ok := tempUserId.(int64)
	if !ok || userId == 0 {
		return
	}
	c.Set(sparkUserId, userId)
	c.Set(sparkLoginStatus, 1)
}

func processResponse(c *gin.Context) {
	Isession, exists := c.Get(sparkSessionName)
	if !exists {
		return
	}
	redisSession, ok := Isession.(session.ISession)
	if !ok {
		return
	}
	if !redisSession.IsModify() {
		return
	}
	err := redisSession.Save()
	if err != nil {
		return
	}
	sessionId := redisSession.Id()
	cookie := &http.Cookie{
		Name:     cookieSessionId,
		Value:    sessionId,
		Path:     "/",
		Domain:   "127.0.0.1:8080",
		MaxAge:   cookieMaxAge,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, cookie)

	return
}

func GetUserId(c *gin.Context) (userId int64, err error) {
	IUserId, exists := c.Get(sparkUserId)
	if !exists {
		err = errors.New("user id not exists")
		return
	}
	userId, ok := IUserId.(int64)
	if !ok {
		err = errors.New("user id convert failed")
		return
	}

	return
}

func IsLogin(c *gin.Context) (login bool, err error) {
	IloginStatus, exists := c.Get(sparkLoginStatus)
	if !exists {
		err = errors.New("login status not exists")
		return
	}
	loginStatus, ok := IloginStatus.(int)
	if !ok {
		err = errors.New("login status convert failed")
		return
	}

	if loginStatus == 1 {
		login = true
	}
	return
}
