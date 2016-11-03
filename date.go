package utils

import(
	"time"
)

func DateNowStringYMD() string {
	return time.Now().Format("2006-01-02")
}

func DateNowNanosecond() int64{
	return time.Now().UnixNano()
}

func DateNowMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)
}

func DateNowSecond() int64 {
	return time.Now().Unix()
}


func DateDurationFrom(t time.Time) time.Duration{
	return time.Now().Sub(t)
}


func DateUnix(date string) int64{
	t, err := time.Parse("2006-01-02", date)
	if err != nil{
		return 0
	}
	return t.Unix()
}

func DateUnixByLocation(date string, location string) int64{
	loc, err := time.LoadLocation(location)
	var t time.Time
	var pErr error
	if err != nil{
		t, pErr = time.Parse("2006-01-02", date)
	}else{
		t, pErr = time.ParseInLocation("2006-01-02", date, loc)
	}
	if pErr != nil{
		return 0
	}
	return t.Unix()
}
