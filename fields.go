package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	types  []string
	fields []field
	custom []field
)

type field struct {
	Key        string   `json:"key"`
	Name       string   `json:"name,omitempty"`
	Size       int      `json:"size,omitempty"`
	Height     string   `json:"height,omitempty"`
	Title      bool     `json:"title,omitempty"`
	Searchable bool     `json:"searchable,omitempty"`
	Required   bool     `json:"required,omitempty"`
	Enum       []string `json:"enum,omitempty"`
}

func parseFields(fields []field) (types []string, err error) {
	for _, i := range fields {
		if i.Key == "type" {
			types = i.Enum
			if i.Required && len(types) == 0 {
				return nil, errors.New("no types")
			}
		}
	}
	return
}

func loadFields() error {
	b, err := os.ReadFile(joinPath(dir(self), "fields.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}
	if types, err = parseFields(fields); err != nil {
		return err
	}
	b, err = os.ReadFile(joinPath(dir(self), "custom.json"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if err == nil {
		if err := json.Unmarshal(b, &custom); err != nil {
			svc.Print(err)
		}
	}
	return nil
}

func updateFields(c *gin.Context) {
	var data []field
	if err := c.BindJSON(&data); err != nil {
		c.String(400, "400")
		return
	}

	res, err := parseFields(data)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	infoMutex.Lock()
	defer infoMutex.Unlock()

	fields = data
	types = res
	b, _ := json.MarshalIndent(data, "", "  ")
	if err := os.WriteFile(joinPath(dir(self), "fields.json"), b, 0644); err != nil {
		svc.Print(err)
		c.String(500, "内部错误")
		return
	}
	c.String(200, "done")
}

func updateCustom(c *gin.Context) {
	var data []field
	if err := c.BindJSON(&data); err != nil {
		c.String(400, "400")
		return
	}

	infoMutex.Lock()
	defer infoMutex.Unlock()

	custom = data
	b, _ := json.MarshalIndent(data, "", "  ")
	if err := os.WriteFile(joinPath(dir(self), "custom.json"), b, 0644); err != nil {
		svc.Print(err)
		c.String(500, "内部错误")
		return
	}
	c.String(200, "done")
}
