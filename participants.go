package main

import (
	"errors"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/txt"
)

var participants []string

func loadParticipants() (err error) {
	participants, err = txt.ReadFile(joinPath(dir(self), "participants.txt"))
	if err != nil {
		return
	}
	if len(participants) == 0 {
		err = errors.New("no participants")
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

	participants = slices.DeleteFunc(data, func(s string) bool { return strings.TrimSpace(s) == "" })
	if err := txt.ExportFile(participants, joinPath(dir(self), "participants.txt")); err != nil {
		svc.Print(err)
		c.String(500, "内部错误")
		return
	}
	c.String(200, "done")
}
