package main

import (
	"Demo/crontab/internal/crontab"
)

var ct *crontab.Crontab

func main() {
	ct = crontab.New()
	ct.Start()
}

//func stop(w http.ResponseWriter, req *http.Request) {
//	req.ParseForm()
//	id := req.Form.Get("id")
//	cronLog := (&models.CronLog{ID: goutils.ToInt(id)}).Get()
//	fmt.Fprintf(w, "cronlog:%+v\n", cronLog)
//}
