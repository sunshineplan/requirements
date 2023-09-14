package main

import (
	"bytes"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/csv"
)

var (
	requirementsList = make(map[ID]requirement)
	lastID           ID

	mu sync.Mutex
)

func get(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	var json []requirement
	for _, v := range requirementsList {
		json = append(json, v)
	}
	slices.SortFunc(json, func(a, b requirement) int { return a.ID.compare(b.ID) })
	c.SetCookie("last", last.String(), 856400*365, "", "", false, true)
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
	svc.Printf("%s %v add %#v", c.ClientIP(), username, data)
	if err := save(); err != nil {
		svc.Println(c.ClientIP(), username, err)
		c.AbortWithStatus(500)
		return
	}
	c.SetCookie("last", last.String(), 856400*365, "", "", false, true)
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
		c.JSON(200, gin.H{"status": 0, "message": "Requirement describe is empty.", "error": 1})
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
	svc.Printf("%s %v edit %#v", c.ClientIP(), username, data.New)
	if err := save(); err != nil {
		svc.Println(c.ClientIP(), username, err)
		c.AbortWithStatus(500)
		return
	}
	c.SetCookie("last", last.String(), 856400*365, "", "", false, true)
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
		svc.Printf("%s %v delete %#v", c.ClientIP(), username, v)
		delete(requirementsList, id)
		if err := save(); err != nil {
			svc.Println(c.ClientIP(), username, err)
			c.AbortWithStatus(500)
			return
		}
		c.SetCookie("last", last.String(), 856400*365, "", "", false, true)
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
	var buf bytes.Buffer
	if err := csv.ExportUTF8(nil, res, &buf); err != nil {
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

func analyze(src []requirement, year, startMonth, endMonth int, isNew bool) (res []summary) {
	for i := startMonth; i <= endMonth; i++ {
		sum := summary{Year: year, Month: i}
		for _, i := range src {
			if (isNew && i.Date.year == sum.Year && i.Date.month == sum.Month) ||
				(!isNew && inRange(i.Date, i.Deadline, sum.Year, sum.Month)) {
				switch i.Type {
				case 1:
					sum.I++
				case 2:
					sum.II++
				case 3:
					sum.III++
				case 4:
					sum.IV++
				case 5:
					sum.V++
				case 6:
					sum.VI++
				case 7:
					sum.VII++
				default:
					svc.Println("unknown requirement type", i)
					continue
				}
				sum.Totle++
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
