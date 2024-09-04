package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/password"
	"github.com/sunshineplan/utils/txt"
)

var (
	infoMutex sync.Mutex
	users     [][2]string
)

func loadUsers() (err error) {
	s, err := txt.ReadFile(joinPath(dir(self), "auth.txt"))
	if err != nil {
		return
	}
	for _, i := range s {
		s := strings.Split(i, ":")
		if len(s) == 2 {
			users = append(users, [2]string{strings.TrimSpace(s[0]), strings.TrimSpace(s[1])})
		}
	}
	if len(users) == 0 {
		err = errors.New("empty user list")
	}
	return
}

func saveUsers() error {
	var s []string
	for _, i := range users {
		s = append(s, i[0]+":"+i[1])
	}
	return txt.ExportFile(s, joinPath(dir(self), "auth.txt"))
}

func getPassword(username string) (string, error) {
	infoMutex.Lock()
	defer infoMutex.Unlock()
	for _, i := range users {
		if i[0] == username {
			return i[1], nil
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
		c.AbortWithStatus(401)
	} else {
		c.Set("username", username)
	}
}

func adminRequired(c *gin.Context) {
	if username := sessions.Default(c).Get("username"); username != "admin" {
		c.AbortWithStatus(403)
	}
}

type info struct {
	username string
	ip       string
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
		if err := password.Compare(info{login.Username, c.ClientIP()}, pwd, login.Password); err != nil {
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

func updateUser(c *gin.Context, add bool) {
	var data struct{ Username, Password string }
	if err := c.BindJSON(&data); err != nil {
		c.String(400, "")
		return
	}
	data.Username = strings.TrimSpace(strings.ToLower(data.Username))

	infoMutex.Lock()
	defer infoMutex.Unlock()

	for i := range users {
		if data.Username == users[i][0] {
			if add {
				c.String(400, "用户已存在")
				return
			} else {
				users[i][1] = data.Password
			}
		}
	}
	if add {
		users = append(users, [2]string{data.Username, data.Password})
	}
	if err := saveUsers(); err != nil {
		svc.Print(err)
		c.String(500, "内部错误")
		return
	}
	c.String(200, "done")
}

func deleteUser(c *gin.Context) {
	var data struct{ Username string }
	if err := c.BindJSON(&data); err != nil {
		c.String(400, "")
		return
	}

	infoMutex.Lock()
	defer infoMutex.Unlock()

	var found bool
	users = slices.DeleteFunc(users, func(s [2]string) bool {
		res := s[0] == data.Username
		if res {
			found = true
		}
		return res
	})

	if found {
		if err := saveUsers(); err != nil {
			svc.Print(err)
			c.String(500, "内部错误")
			return
		}
		c.String(200, "done")
	} else {
		c.String(400, "用户不存在")
	}
}
