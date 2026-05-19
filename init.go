package main

import (
	"errors"
	"io/fs"

	"github.com/sunshineplan/utils/csv"
)

func initSrv() error {
	var res []requirement
	if err := csv.DecodeFile(joinPath(dir(self), "requirements.csv"), &res); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, i := range res {
		requirementsList[i.ID] = i
		if i.ID.compare(lastID) > 0 {
			lastID = i.ID
		}
	}
	if err := setLast(); err != nil {
		svc.Fatal(err)
	}
	return nil
}

func test() error {
	return initSrv()
}
