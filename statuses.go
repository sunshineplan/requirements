package main

import (
	"errors"
	"io/fs"
	"strconv"
	"strings"
	"unicode"

	"github.com/sunshineplan/utils/txt"
)

var statuses []status

type status struct {
	Value  string `json:"value"`
	Closed bool   `json:"closed"`
}

func parseStatus(s string) status {
	res := strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || r == ':'
	})
	value := res[0]
	var closed bool
	if len(res) > 1 {
		b, _ := strconv.ParseBool(res[1])
		closed = !b
	}
	if value == *doneValue {
		closed = true
	}
	return status{value, closed}
}

func loadStatuses() error {
	s, err := txt.ReadFile(joinPath(dir(self), "statuses.txt"))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	for _, i := range s {
		if s := strings.TrimSpace(i); s != "" && s != ":" {
			statuses = append(statuses, parseStatus(s))
		}
	}
	return nil
}
