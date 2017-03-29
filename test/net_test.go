package test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/heqzha/goutils/net"
)

func TestHTTPGetEncodeQuery(t *testing.T) {
	u := "https%3A%2F%2Fwww.google.com.hk%2F%23safe%3Dstrict%26q%3D%25E6%25B5%258B%25E8%25AF%2595%26*"
	unescape, _ := url.QueryUnescape(u)
	fmt.Println(unescape)
	o, _ := url.Parse(unescape)
	q := o.Query()
	o.RawQuery = q.Encode() //Important: Encode query, otherwise http.NewRequest will return 400 error
	fmt.Println(o.String())
	r, err := net.HTTPGet(o.String(), nil, nil)
	fmt.Println(err, string(r))
}
