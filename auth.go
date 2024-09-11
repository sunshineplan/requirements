package main

import (
	"cmp"
	"errors"
	"fmt"
	"maps"
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
	users     = make(map[string]string)
)

func loadUsers() (err error) {
	s, err := txt.ReadFile(joinPath(dir(self), "auth.txt"))
	if err != nil {
		return
	}
	var admin bool
	for _, i := range s {
		s := strings.Split(i, ":")
		if len(s) == 2 {
			username, password := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
			if username == "" || password == "" {
				svc.Printf("bad user info: %q", i)
				continue
			}
			if _, ok := users[username]; ok {
				continue
			}
			if username == "admin" {
				admin = true
			}
			users[username] = password
		}
	}
	if !admin {
		users["admin"] = "123456"
	}
	return
}

func usernames() (usernames []string) {
	for username := range maps.Keys(users) {
		usernames = append(usernames, username)
	}
	slices.SortFunc(usernames, func(a, b string) int {
		if a == "admin" {
			return -1
		}
		if b == "admin" {
			return 1
		}
		return cmp.Compare(a, b)
	})
	return
}

func saveUsers() error {
	var s []string
	for _, i := range usernames() {
		s = append(s, i+":"+users[i])
	}
	return txt.ExportFile(s, joinPath(dir(self), "auth.txt"))
}

func getPassword(username string) (string, error) {
	infoMutex.Lock()
	defer infoMutex.Unlock()
	if password, ok := users[username]; ok {
		return password, nil
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

	if _, ok := users[data.Username]; ok {
		if add {
			c.String(400, "用户已存在")
			return
		} else {
			users[data.Username] = data.Password
		}
	} else if add {
		users[data.Username] = data.Password
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
		c.Status(400)
		return
	}
	if data.Username == "admin" {
		c.String(400, "禁止删除admin")
		return
	}

	infoMutex.Lock()
	defer infoMutex.Unlock()

	if _, ok := users[data.Username]; ok {
		delete(users, data.Username)
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
