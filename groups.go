package main

import (
	"errors"
	"io/fs"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/txt"
)

var groups []string

func loadGroups() (err error) {
	if groups, err = txt.ReadFile(joinPath(dir(self), "groups.txt")); groups == nil {
		groups = []string{}
	}
	if errors.Is(err, fs.ErrNotExist) {
		err = nil
	}
	return
}

func updateGroups(c *gin.Context) {
	var data []string
	if err := c.BindJSON(&data); err != nil {
		c.String(400, "")
		return
	}

	infoMutex.Lock()
	defer infoMutex.Unlock()

	if groups = slices.DeleteFunc(data, func(s string) bool {
		return strings.TrimSpace(s) == ""
	}); groups == nil {
		groups = []string{}
	}
	if err := txt.ExportFile(groups, joinPath(dir(self), "groups.txt")); err != nil {
		svc.Print(err)
		c.String(500, "内部错误")
		return
	}
	c.String(200, "done")
}
