package main

import (
	"crypto/rand"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var dist embed.FS

func run() error {
	if err := loadUsers(); err != nil {
		return err
	}
	if err := loadFields(); err != nil {
		return err
	}

	router := gin.Default()
	router.TrustedPlatform = "X-Real-IP"
	server.Handler = router

	if err := initSrv(); err != nil {
		return err
	}

	secret := make([]byte, 16)
	if _, err := rand.Read(secret); err != nil {
		return err
	}
	router.Use(sessions.Sessions("session", cookie.NewStore(secret)))

	assets, err := fs.Sub(dist, "dist/assets")
	if err != nil {
		return err
	}
	router.StaticFS("/assets", http.FS(assets))
	router.StaticFileFS("favicon.ico", "dist/favicon.ico", http.FS(dist))
	router.SetHTMLTemplate(template.Must(template.New("").ParseFS(dist, "dist/*.html")))

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.GET("/info", func(c *gin.Context) {
		obj := map[string]any{"brand": *brand}
		user, _ := getUser(c)
		if user == "" {
			c.JSON(200, obj)
			return
		}
		infoMutex.Lock()
		defer infoMutex.Unlock()
		obj["username"] = user
		obj["fields"] = fields
		if len(custom) > 0 {
			obj["custom"] = custom
		}
		obj["done"] = *doneValue
		if user == "admin" {
			obj["users"] = usernames()
		}
		if last.Equal(c) {
			c.JSON(200, obj)
		} else {
			c.SetCookie("last", last.String(), 856400*365, "", "", false, false)
			c.Status(409)
		}
	})
	router.GET("/poll", authRequired, func(c *gin.Context) {
		time.Sleep(*poll)
		infoMutex.Lock()
		defer infoMutex.Unlock()
		c.String(200, last.String())
	})

	auth := router.Group("/")
	auth.POST("/login", login)
	auth.POST("/logout", authRequired, func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		if err := session.Save(); err != nil {
			svc.Print(err)
			c.String(500, "")
			return
		}
		c.String(200, "bye")
	})

	base := router.Group("/", authRequired)
	base.GET("/get", get)
	base.GET("/statistics", statistics)
	base.POST("/add", add)
	base.POST("/edit", edit)
	base.POST("/done", done)

	admin := base.Group("/", adminRequired)
	admin.POST("/delete/:id", del)
	admin.POST("/fields", updateFields)
	admin.POST("/custom", updateCustom)
	admin.POST("/addUser", func(c *gin.Context) { updateUser(c, true) })
	admin.POST("/chgpwd", func(c *gin.Context) { updateUser(c, false) })
	admin.POST("/deleteUser", deleteUser)

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})

	return server.Run()
}
