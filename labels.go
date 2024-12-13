package main

import (
	"errors"
	"io/fs"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/txt"
)

var labels []string

func loadLabels() (err error) {
	if labels, err = txt.ReadFile(joinPath(dir(self), "labels.txt")); labels == nil {
		labels = []string{}
	}
	if errors.Is(err, fs.ErrNotExist) {
		err = nil
	}
	return
}

func updateLabels(c *gin.Context) {
	var data []string
	if err := c.BindJSON(&data); err != nil {
		c.String(400, "")
		return
	}

	infoMutex.Lock()
	defer infoMutex.Unlock()

	if labels = slices.DeleteFunc(data, func(s string) bool {
		return strings.TrimSpace(s) == ""
	}); labels == nil {
		labels = []string{}
	}
	if err := txt.ExportFile(labels, joinPath(dir(self), "labels.txt")); err != nil {
		svc.Print(err)
		c.String(500, "内部错误")
		return
	}
	c.String(200, "done")
}
