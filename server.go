package main

import (
	"crypto/rand"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/log"
)

func run() error {
	if *logPath != "" {
		svc.Logger = log.New(*logPath, "", log.LstdFlags)
		gin.DefaultWriter = svc.Logger
		gin.DefaultErrorWriter = svc.Logger
	}

	if err := loadUsers(); err != nil {
		return err
	}
	if err := loadParticipants(); err != nil {
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

	router.StaticFS("/assets", http.Dir(joinPath(dir(self), "dist/assets")))
	router.StaticFile("favicon.ico", joinPath(dir(self), "dist/favicon.ico"))
	router.LoadHTMLFiles(joinPath(dir(self), "dist/index.html"))

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.GET("/info", func(c *gin.Context) {
		user, _ := getUser(c)
		if user == "" {
			c.JSON(200, struct{}{})
			return
		}
		infoMutex.Lock()
		defer infoMutex.Unlock()
		obj := map[string]any{"username": user, "participants": participants}
		if user == "admin" {
			var s []string
			for _, i := range users {
				s = append(s, i[0])
			}
			obj["users"] = s
		}
		if last.Equal(c) {
			c.JSON(200, obj)
		} else {
			c.SetCookie("last", last.String(), 856400*365, "", "", false, true)
			c.AbortWithStatus(409)
		}
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
	base.POST("/delete/:id", del)

	admin := base.Group("/", adminRequired)
	admin.POST("/participants", updateParticipants)
	admin.POST("/addUser", func(c *gin.Context) { updateUser(c, true) })
	admin.POST("/chgpwd", func(c *gin.Context) { updateUser(c, false) })
	admin.POST("/deleteUser", deleteUser)

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})

	return server.Run()
}
