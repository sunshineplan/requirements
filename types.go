package main

import (
	"errors"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/txt"
)

var types []string

func loadTypes() (err error) {
	types, err = txt.ReadFile(joinPath(dir(self), "types.txt"))
	if err != nil {
		return
	}
	if len(types) == 0 {
		err = errors.New("no types")
	}
	return
}

func updateTypes(c *gin.Context) {
	var data []string
	if err := c.BindJSON(&data); err != nil {
		c.String(400, "")
		return
	}

	infoMutex.Lock()
	defer infoMutex.Unlock()

	types = slices.DeleteFunc(data, func(s string) bool { return strings.TrimSpace(s) == "" })
	if err := txt.ExportFile(types, joinPath(dir(self), "types.txt")); err != nil {
		svc.Print(err)
		c.String(500, "内部错误")
		return
	}
	c.String(200, "done")
}
