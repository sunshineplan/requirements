package main

import (
	"cmp"
	"encoding"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	_ encoding.TextUnmarshaler = new(ID)
	_ encoding.TextUnmarshaler = new(Date)
	_ encoding.TextMarshaler   = ID{}
	_ encoding.TextMarshaler   = Date{}
)

type ID struct {
	year   int
	serial int
}

func (id *ID) UnmarshalText(text []byte) error {
	if res := regexp.MustCompile(`^(\d{4})(\d{4})$`).FindStringSubmatch(string(text)); len(res) != 3 {
		return fmt.Errorf("unknown id: %s", text)
	} else if year, err := time.Parse("2006", res[1]); err != nil {
		return fmt.Errorf("unknown id: %s", text)
	} else {
		serial, _ := strconv.Atoi(res[2])
		*id = ID{year.Year(), serial}
	}
	return nil
}

func (id ID) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%d%04d", id.year, id.serial)), nil
}

func newID(last ID) ID {
	year := time.Now().Year()
	if last.year == 0 && last.serial == 0 {
		return ID{year, 1}
	}
	switch cmp.Compare(year, last.year) {
	case 1:
		return ID{year, 1}
	case 0:
		return ID{year, last.serial + 1}
	default:
		return ID{last.year, last.serial + 1}
	}
}

func parseID(s string) (id ID, err error) {
	err = id.UnmarshalText([]byte(s))
	return
}

func (x ID) compare(y ID) int {
	if n := cmp.Compare(x.year, y.year); n == 0 {
		return cmp.Compare(x.serial, y.serial)
	} else {
		return n
	}
}

type Date struct {
	year, month, day int
}

func now() Date {
	y, m, d := time.Now().Date()
	return Date{y, int(m), d}
}

func (d *Date) UnmarshalText(text []byte) error {
	s := strings.Map(func(r rune) rune {
		if r == '-' || r == '/' || unicode.IsSpace(r) {
			return -1
		}
		return r
	}, string(text))
	if t, err := time.Parse("20060102", s); err == nil {
		*d = Date{t.Year(), int(t.Month()), t.Day()}
	}
	return nil
}

func (d Date) MarshalText() ([]byte, error) {
	if d.isZero() {
		return nil, nil
	}
	return []byte(d.String()), nil
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

func (date Date) isZero() bool {
	return date.year == 0 && date.month == 0 && date.day == 0
}
