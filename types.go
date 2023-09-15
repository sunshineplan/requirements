package main

import (
	"cmp"
	"encoding"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	_ encoding.TextUnmarshaler = new(ID)
	_ encoding.TextUnmarshaler = new(Type)
	_ encoding.TextUnmarshaler = new(Date)
	_ encoding.TextMarshaler   = ID{}
	_ encoding.TextMarshaler   = Type(0)
	_ encoding.TextMarshaler   = Date{}
)

type requirement struct {
	ID            ID     `csv:"编号" json:"id"`
	Type          Type   `csv:"类型" json:"type"`
	Desc          string `csv:"描述" json:"desc"`
	Date          Date   `csv:"提请日期" json:"date"`
	Deadline      Date   `csv:"期限日期" json:"deadline"`
	Submitter     string `csv:"提交人" json:"submitter"`
	Recipient     string `csv:"承接人" json:"recipient"`
	Acceptor      string `csv:"受理人" json:"acceptor"`
	Status        string `csv:"状态" json:"status"`
	Note          string `csv:"备注" json:"note"`
	Participating string `csv:"参与班组" json:"participating"`
}

func (r requirement) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

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
	switch year := time.Now().Year(); cmp.Compare(year, last.year) {
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

type Type int

var typeList = []string{
	"内容策划",
	"宣传推广",
	"用户培训",
	"宣传品相关",
	"平台相关",
	"中心业务",
	"馆所业务",
}

func (t *Type) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))
	if res := regexp.MustCompile(`^(\d)、`).FindStringSubmatch(s); len(res) != 0 {
		v, _ := strconv.Atoi(res[1])
		*t = Type(v)
	} else if i := slices.Index(typeList, s); i != -1 {
		*t = Type(i + 1)
	}
	return nil
}

func (t Type) MarshalText() ([]byte, error) {
	if t == 0 {
		return []byte("未知"), nil
	}
	return []byte(typeList[t-1]), nil
}

type Date struct {
	year, month, day int
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
	return []byte(fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)), nil
}

func (date Date) isZero() bool {
	return date.year == 0 && date.month == 0 && date.day == 0
}

func inRange(date, deadline Date, year, month int) bool {
	if date.isZero() && deadline.isZero() {
		return true
	}
	if date.isZero() {
		if year == deadline.year {
			return month <= deadline.month
		}
		return year <= deadline.year
	}
	if deadline.isZero() {
		if year == date.year {
			return month >= date.month
		}
		return year >= date.year
	}
	d1 := fmt.Sprintf("%d%02d", date.year, date.month)
	d2 := fmt.Sprintf("%d%02d", deadline.year, deadline.month)
	ym := fmt.Sprintf("%d%02d", year, month)
	return d1 <= ym && ym <= d2
}

type summary struct {
	Year  int `csv:"年份"`
	Month int `csv:"月"`
	I     int `csv:"内容策划"`
	II    int `csv:"宣传推广"`
	III   int `csv:"用户培训"`
	IV    int `csv:"宣传品相关"`
	V     int `csv:"平台相关"`
	VI    int `csv:"中心业务"`
	VII   int `csv:"馆所业务"`
	Totle int `csv:"总计"`
}
