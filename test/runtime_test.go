package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/heqzha/goutils/date"
	"github.com/heqzha/goutils/net"
	"github.com/heqzha/goutils/runtime"
)

func HTTPGet(url string) (string, error) {
	resp, err := net.HTTPGet(url, nil, nil)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func TimeCost() {
	now := time.Now()
	defer fmt.Printf("%s cost %s\n", "HTTPGet direct", date.DateDurationFrom(now))
	HTTPGet("http://www.baidu.com")
}

func TestGetTimeCost(t *testing.T) {
	TimeCost()
	out := runtime.PrintTimeCost(HTTPGet, "http://www.baidu.com")
	if out[1] != nil {
		fmt.Println(out[1].(error))
	}
}
