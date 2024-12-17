package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
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

var (
	_ json.Marshaler   = requirement{}
	_ json.Unmarshaler = new(requirement)
)

type requirement struct {
	ID        ID     `csv:"id" json:"id"`
	Type      string `csv:"type" json:"type"`
	Title     string `csv:"title" json:"title"`
	Date      Date   `csv:"date" json:"date"`
	Deadline  Date   `csv:"deadline" json:"deadline"`
	Done      Date   `csv:"done" json:"done"`
	Submitter string `csv:"submitter" json:"submitter"`
	Recipient string `csv:"recipient" json:"recipient"`
	Acceptor  string `csv:"acceptor" json:"acceptor"`
	Status    string `csv:"status" json:"status"`
	Label     string `csv:"label" json:"label"`
	Note      string `csv:"note" json:"note"`

	customFields map[string]string
}

func (r requirement) MarshalJSON() ([]byte, error) {
	v := reflect.ValueOf(r)
	t := v.Type()
	m := make(map[string]string)
	for _, field := range reflect.VisibleFields(t) {
		if !field.IsExported() {
			continue
		}
		m[strings.ToLower(field.Name)] = fmt.Sprint(v.FieldByIndex(field.Index))
	}
	maps.Insert(m, maps.All(r.customFields))
	return json.Marshal(m)
}

func (r *requirement) UnmarshalJSON(b []byte) error {
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	requirement := requirement{customFields: make(map[string]string)}
	for k, v := range m {
		switch strings.ToLower(k) {
		case "id":
			id, err := parseID(v)
			if err != nil {
				return err
			}
			requirement.ID = id
		case "type":
			requirement.Type = v
		case "title":
			requirement.Title = v
		case "date":
			date, err := parseDate(v)
			if err != nil {
				return err
			}
			requirement.Date = date
		case "deadline":
			deadline, err := parseDate(v)
			if err != nil {
				return err
			}
			requirement.Deadline = deadline
		case "done":
			done, err := parseDate(v)
			if err != nil {
				return err
			}
			requirement.Done = done
		case "submitter":
			requirement.Submitter = v
		case "recipient":
			requirement.Recipient = v
		case "acceptor":
			requirement.Acceptor = v
		case "status":
			requirement.Status = v
		case "label":
			requirement.Label = v
		case "note":
			requirement.Note = v
		default:
			requirement.customFields[k] = v
		}
	}
	*r = requirement
	return nil
}

func (r requirement) IsEqual(u requirement) bool {
	a, b := reflect.ValueOf(r), reflect.ValueOf(u)
	t := a.Type()
	for _, field := range reflect.VisibleFields(t) {
		if !field.IsExported() {
			continue
		}
		if !a.FieldByIndex(field.Index).Equal(b.FieldByIndex(field.Index)) {
			return false
		}
	}
	return maps.Equal(r.customFields, u.customFields)
}

func (r requirement) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func get(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	json := []requirement{}
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

	if data.Title == "" {
		c.JSON(200, gin.H{"status": 0, "message": "标题为空。", "error": 1})
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
		fmt.Sprintf("%s新增了一项记录-%s", username, time.Now().Format("20060102 15:04")),
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

	if data.New.Title == "" {
		c.JSON(200, gin.H{"status": 0, "message": "标题为空。", "error": 1})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if v := requirementsList[data.Old.ID]; !v.IsEqual(data.Old) {
		c.AbortWithStatus(409)
		return
	}
	requirementsList[data.New.ID] = data.New
	obj := gin.H{"status": 1}
	if !last.Equal(c) {
		obj["reload"] = 1
	}
	svc.Printf("%s %v edit %s", c.ClientIP(), username, data.New)
	go sendMail(
		fmt.Sprintf("%s编辑了一项记录-%s", username, time.Now().Format("20060102 15:04")),
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

func done(c *gin.Context) {
	username, _ := c.Get("username")
	var data requirement
	if err := c.BindJSON(&data); err != nil {
		svc.Println(c.ClientIP(), username, err)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if v := requirementsList[data.ID]; !v.IsEqual(data) {
		c.AbortWithStatus(409)
		return
	}
	date, err := parseDate(c.Query("date"))
	if err != nil {
		date = now()
	}
	obj := gin.H{"status": 1}
	if !last.Equal(c) {
		obj["reload"] = 1
	}
	if data.Status == *doneValue {
		c.JSON(200, obj)
		return
	} else {
		data.Status = *doneValue
		data.Done = date
		obj["value"] = data.Status
		obj["done"] = data.Done
	}
	requirementsList[data.ID] = data
	svc.Printf("%s %v done %s", c.ClientIP(), username, data)
	go sendMail(
		fmt.Sprintf("%s完成了一项记录-%s", username, time.Now().Format("20060102 15:04")),
		fmt.Sprintf("完成内容:\n%s\n\nIP: %s", data, c.ClientIP()),
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
			fmt.Sprintf("%s删除了一项记录-%s", username, time.Now().Format("20060102 15:04")),
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
	var s []string
	for _, i := range types {
		s = append(s, i)
		s = append(s, i+*doneValue)
	}
	fieldnames := append(append([]string{"年份", "月"}, s...), "总计", "总计"+*doneValue)
	var data []map[string]int
	for _, i := range res {
		m := i.Types
		m["年份"] = i.Year
		m["月"] = i.Month
		m["总计"] = i.Total
		m["总计"+*doneValue] = i.TotalDone
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

func getFields() []string {
	fieldnames := []string{
		"id",
		"type",
		"title",
		"date",
		"deadline",
		"done",
		"submitter",
		"recipient",
		"acceptor",
		"status",
		"label",
		"note",
	}
	for _, field := range custom {
		fieldnames = append(fieldnames, field.Key)
	}
	return fieldnames
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
	if err := csv.ExportUTF8(getFields(), rows, f); err != nil {
		return err
	}
	return setLast()
}

func backup() {
	mu.Lock()
	defer mu.Unlock()
	if len(requirementsList) == 0 {
		return
	}
	sendMail(
		fmt.Sprintf("数据备份-%s", time.Now().Format("20060102")),
		fmt.Sprintf("备份时间: %s", time.Now()),
		[]*mail.Attachment{{Path: joinPath(dir(self), "requirements.csv")}},
	)
}

type summary struct {
	Year      int
	Month     int
	Types     map[string]int
	Total     int
	TotalDone int
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
			if i.Done.year == sum.Year && i.Done.month == sum.Month {
				sum.Types[i.Type+*doneValue]++
				sum.TotalDone++
			}
		}
		for _, i := range types {
			if _, ok := sum.Types[i]; !ok {
				sum.Types[i] = 0
			}
			if _, ok := sum.Types[i+*doneValue]; !ok {
				sum.Types[i+*doneValue] = 0
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
