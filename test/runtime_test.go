package test

import (
	"fmt"
	"testing"

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

func TestGetTimeCost(t *testing.T) {
	out := runtime.PrintTimeCost(HTTPGet, "http://www.baidu.com")
	if out[1] != nil {
		fmt.Println(out[1].(error))
	}
	fmt.Println(len(out[0].(string)))
}
