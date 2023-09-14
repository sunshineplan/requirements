package main

import "github.com/sunshineplan/utils/csv"

func initSrv() error {
	//var metadata any
	//if err := retry.Do(func() error {
	//	return meta.Get("requirements", &metadata)
	//}, 3, 20); err != nil {
	//	return err
	//}

	//if metadata == nil {
	//	return errors.New("no permission")
	//}

	var res []requirement
	if err := csv.DecodeFile(joinPath(dir(self), "requirements.csv"), &res); err != nil {
		return err
	}
	for _, i := range res {
		requirementsList[i.ID] = i
		if i.ID.compare(lastID) > 0 {
			lastID = i.ID
		}
	}
	return nil
}

func test() error {
	return initSrv()
}
