package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/csv"
	"github.com/sunshineplan/utils/mail"
)

var (
	requirementsList = make(map[ID]requirement)
	lastID           ID

	mu sync.Mutex
)

type requirement struct {
	ID            ID     `csv:"编号" json:"id"`
	Type          string `csv:"类型" json:"type"`
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

func get(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	var json []requirement
	for _, v := range requirementsList {
		json = append(json, v)
	}
	slices.SortFunc(json, func(a, b requirement) int { return a.ID.compare(b.ID) })
	c.SetCookie("last", last.String(), 856400*365, "", "", false, false)
	c.JSON(200, json)
}

func add(c *gin.Context) {
	username, _ := c.Get("username")
	var data requirement
	if err := c.BindJSON(&data); err != nil {
		svc.Println(c.ClientIP(), username, err)
		return
	}

	if data.Desc == "" {
		c.JSON(200, gin.H{"status": 0, "message": "Requirement describe is empty.", "error": 1})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	lastID = newID(lastID)
	data.ID = lastID
	requirementsList[data.ID] = data
	obj := gin.H{"status": 1, "id": data.ID}
	if !last.Equal(c) {
		obj["reload"] = 1
	}
	svc.Printf("%s %v add %s", c.ClientIP(), username, data)
	go sendMail(
		fmt.Sprintf("[业务系统]%s新增了一项业务-%s", username, time.Now().Format("20060102 15:04")),
		fmt.Sprintf("%s\n\nIP: %s", data, c.ClientIP()),
		nil,
	)
	if err := save(); err != nil {
		svc.Println(c.ClientIP(), username, err)
		c.AbortWithStatus(500)
		return
	}
	c.SetCookie("last", last.String(), 856400*365, "", "", false, false)
	c.JSON(200, obj)
}

func edit(c *gin.Context) {
	username, _ := c.Get("username")
	var data struct {
		Old, New requirement
	}
	if err := c.BindJSON(&data); err != nil {
		svc.Println(c.ClientIP(), username, err)
		return
	}

	if data.New.Desc == "" {
		c.JSON(200, gin.H{"status": 0, "message": "业务描述为空。", "error": 1})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if v := requirementsList[data.Old.ID]; v != data.Old {
		c.AbortWithStatus(409)
	}
	requirementsList[data.New.ID] = data.New
	obj := gin.H{"status": 1}
	if !last.Equal(c) {
		obj["reload"] = 1
	}
	svc.Printf("%s %v edit %s", c.ClientIP(), username, data.New)
	go sendMail(
		fmt.Sprintf("[业务系统]%s编辑了一项业务-%s", username, time.Now().Format("20060102 15:04")),
		fmt.Sprintf("原始内容:\n%s\n\n修改内容:\n%s\n\nIP: %s", data.Old, data.New, c.ClientIP()),
		nil,
	)
	if err := save(); err != nil {
		svc.Println(c.ClientIP(), username, err)
		c.AbortWithStatus(500)
		return
	}
	c.SetCookie("last", last.String(), 856400*365, "", "", false, false)
	c.JSON(200, obj)
}

func del(c *gin.Context) {
	username, _ := c.Get("username")
	id, err := parseID(c.Param("id"))
	if err != nil {
		svc.Println(c.ClientIP(), username, err)
		c.AbortWithStatus(400)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	obj := gin.H{"status": 1}
	if !last.Equal(c) {
		obj["reload"] = 1
	}
	if v, ok := requirementsList[id]; ok {
		svc.Printf("%s %v delete %s", c.ClientIP(), username, v)
		go sendMail(
			fmt.Sprintf("[业务系统]%s删除了一项业务-%s", username, time.Now().Format("20060102 15:04")),
			fmt.Sprintf("%s\n\nIP: %s", v, c.ClientIP()),
			nil,
		)
		delete(requirementsList, id)
		if err := save(); err != nil {
			svc.Println(c.ClientIP(), username, err)
			c.AbortWithStatus(500)
			return
		}
		c.SetCookie("last", last.String(), 856400*365, "", "", false, false)
	}
	c.JSON(200, obj)
}

func statistics(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	username, _ := c.Get("username")
	var src []requirement
	if err := csv.DecodeFile(joinPath(dir(self), "requirements.csv"), &src); err != nil {
		svc.Println(c.ClientIP(), username, err)
		c.AbortWithStatus(500)
		return
	}
	var isNew bool
	if ok, _ := strconv.ParseBool(c.Query("isNew")); ok {
		isNew = true
	}
	res := analyzeFull(src, 2022, 8, isNew)
	fieldnames := append(append([]string{"年份", "月"}, types...), "总计")
	var data []map[string]int
	for _, i := range res {
		m := i.Types
		m["年份"] = i.Year
		m["月"] = i.Month
		m["总计"] = i.Total
		data = append(data, m)
	}
	var buf bytes.Buffer
	if err := csv.ExportUTF8(fieldnames, data, &buf); err != nil {
		svc.Println(c.ClientIP(), username, err)
		c.AbortWithStatus(500)
		return
	}
	c.Data(200, "text/csv", buf.Bytes())
}

func save() error {
	var rows []requirement
	for _, v := range requirementsList {
		rows = append(rows, v)
	}
	slices.SortFunc(rows, func(a, b requirement) int { return a.ID.compare(b.ID) })
	f, err := os.Create(joinPath(dir(self), "requirements.csv"))
	if err != nil {
		return err
	}
	defer f.Close()
	if err := csv.ExportUTF8(nil, rows, f); err != nil {
		return err
	}
	return setLast()
}

func backup() {
	mu.Lock()
	defer mu.Unlock()
	sendMail(
		fmt.Sprintf("[业务系统]数据备份-%s", time.Now().Format("20060102")),
		fmt.Sprintf("备份时间: %s", time.Now()),
		[]*mail.Attachment{{Path: joinPath(dir(self), "requirements.csv")}},
	)
}

type summary struct {
	Year  int
	Month int
	Types map[string]int
	Total int
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

func analyze(src []requirement, year, startMonth, endMonth int, isNew bool) (res []summary) {
	for i := startMonth; i <= endMonth; i++ {
		sum := summary{Year: year, Month: i, Types: make(map[string]int)}
		for _, i := range src {
			if (isNew && i.Date.year == sum.Year && i.Date.month == sum.Month) ||
				(!isNew && inRange(i.Date, i.Deadline, sum.Year, sum.Month)) {
				sum.Types[i.Type]++
				sum.Total++
			}
		}
		for _, i := range types {
			if _, ok := sum.Types[i]; !ok {
				sum.Types[i] = 0
			}
		}
		res = append(res, sum)
	}
	return
}

func analyzeFull(src []requirement, year, month int, isNew bool) (res []summary) {
	t, now := ym(year, month), time.Now()
	if t = ym(t.Year()+1, 1); t.Before(now) {
		res = append(res, analyze(src, year, month, 12, isNew)...)
		for t = ym(t.Year()+1, 1); t.Before(now); t = ym(t.Year()+1, 1) {
			res = append(res, analyze(src, t.Year()-1, 1, 12, isNew)...)
		}
	}
	if t = ym(t.Year()-1, 1); t.Before(now) {
		if year != now.Year() {
			month = 1
		}
		res = append(res, analyze(src, t.Year(), month, int(now.Month()), isNew)...)
	}
	return
}

func ym(year, month int) time.Time {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
}
