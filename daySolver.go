package main

import (
	ptime "github.com/yaa110/go-persian-calendar"
	"time"
)

func progress(start int) int {
	maxDays := 365
	currentTime := time.Now()
	var t = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second(), 0, ptime.Iran())
	pt := ptime.New(t)
	sum := 0
	if start == pt.YearDay() {
		sum = 0
	} else if pt.YearDay() < start {
		sum = (start - maxDays) + pt.YearDay()
	} else if start < currentTime.YearDay() {
		sum = pt.YearDay() - start
	}
	return sum
}

//public int solver() {
//
//if (cDay == day){
//daySum =0;
//Log.i(DAY_LOG , "day is 1 : " +daySum);
//return daySum;
//}else if (cDay<day){
//daySum = Math.abs(day-allDays) + cDay;
//Log.i(DAY_LOG , "day is 2 : " +daySum);
//return daySum;
//}else if (cDay>day){
//daySum = Math.abs(cDay - day);
//Log.i(DAY_LOG , "day is 3 : " +daySum);
//return daySum;
//}
//Log.i(DAY_LOG , "day is : finish " + daySum);
//return daySum;
//}
