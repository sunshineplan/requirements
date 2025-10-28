package main

import (
	"errors"
	"io/fs"
	"time"

	"github.com/sunshineplan/utils/csv"
	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/utils/retry"
)

func initSrv() error {
	var data struct {
		Dialer mail.Dialer
	}
	if err := retry.Do(func() error {
		return meta.Get("requirements", &data)
	}, 3, 20*time.Second); err != nil {
		return err
	}
	if data.Dialer.Server == "" {
		return errors.New("no permission")
	}
	dialer = data.Dialer

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
