package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/password"
	"github.com/sunshineplan/utils/txt"
)

var authMutex sync.Mutex

func getPassword(username string) (string, error) {
	authMutex.Lock()
	defer authMutex.Unlock()
	s, err := txt.ReadFile(joinPath(dir(self), "auth.txt"))
	if err != nil {
		return "", err
	}
	for _, i := range s {
		s := strings.Split(i, ":")
		if len(s) == 2 && strings.TrimSpace(s[0]) == username {
			return strings.TrimSpace(s[1]), nil
		}
	}
	return "", errors.New("user not found")
}

func getUser(c *gin.Context) (user string, ok bool) {
	username := sessions.Default(c).Get("username")
	if username == nil {
		return
	}
	return username.(string), true
}

func authRequired(c *gin.Context) {
	if username := sessions.Default(c).Get("username"); username == nil {
		c.Abort()
		c.Redirect(302, "/")
	} else {
		c.Set("username", username)
	}
}

func login(c *gin.Context) {
	var login struct {
		Username, Password string
		Rememberme         bool
	}
	if err := c.BindJSON(&login); err != nil {
		c.String(400, "")
		return
	}
	login.Username = strings.ToLower(login.Username)

	if password.IsMaxAttempts(c.ClientIP() + login.Username) {
		c.JSON(200, gin.H{"status": 0, "message": fmt.Sprintf("达到最大重试次数(%d)", 5)})
		return
	}

	var message string
	if pwd, err := getPassword(login.Username); err != nil {
		svc.Print(err)
		message = "无效用户名"
	} else {
		if _, err := password.Compare(c.ClientIP()+login.Username, pwd, login.Password, false); err != nil {
			if errors.Is(err, password.ErrIncorrectPassword) {
				message = err.Error()
			} else {
				svc.Print(err)
				c.String(500, "内部错误")
				return
			}
		}

		if message == "" {
			session := sessions.Default(c)
			session.Clear()
			session.Set("username", login.Username)

			if login.Rememberme {
				session.Options(sessions.Options{HttpOnly: true, MaxAge: 856400 * 365})
			} else {
				session.Options(sessions.Options{HttpOnly: true})
			}

			if err := session.Save(); err != nil {
				svc.Print(err)
				c.String(500, "内部错误")
				return
			}

			c.JSON(200, gin.H{"status": 1})
			return
		}
	}

	c.JSON(200, gin.H{"status": 0, "message": message})
}
