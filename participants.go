package main

import (
	"errors"
	"io/fs"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/txt"
)

var participants []string

func loadParticipants() (err error) {
	if participants, err = txt.ReadFile(joinPath(dir(self), "participants.txt")); participants == nil {
		participants = []string{}
	}
	if errors.Is(err, fs.ErrNotExist) {
		err = nil
	}
	return
}

func updateParticipants(c *gin.Context) {
	var data []string
	if err := c.BindJSON(&data); err != nil {
		c.String(400, "")
		return
	}

	infoMutex.Lock()
	defer infoMutex.Unlock()

	if participants = slices.DeleteFunc(data, func(s string) bool {
		return strings.TrimSpace(s) == ""
	}); participants == nil {
		participants = []string{}
	}
	if err := txt.ExportFile(participants, joinPath(dir(self), "participants.txt")); err != nil {
		svc.Print(err)
		c.String(500, "内部错误")
		return
	}
	c.String(200, "done")
}
