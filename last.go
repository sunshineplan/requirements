package main

import (
	"cmp"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var last Last

type Last time.Time

func setLast() error {
	stat, err := os.Stat(joinPath(dir(self), "requirements.csv"))
	if err != nil {
		return err
	}
	last = Last(stat.ModTime())
	return nil
}

func (t Last) String() string {
	return strconv.FormatInt(time.Time(t).UnixMilli(), 10)
}

func (t Last) Equal(c *gin.Context) bool {
	last, _ := c.Cookie("last")
	return cmp.Compare(t.String(), last) == 0
}
