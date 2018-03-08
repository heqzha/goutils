package net

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func HTTPGet(url string, headers map[string]string, cookies []*http.Cookie) ([]byte, error) {
	return HTTPGetWithTimeout(url, headers, cookies, time.Duration(0))
}

func CustomRequest(method string, url string, bodyData []byte) ([]byte, error) {
	b := bytes.NewReader(bodyData)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "plain/text")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Failed to call [%s], status code: %d", url, resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func HTTPGetWithTimeout(url string, headers map[string]string, cookies []*http.Cookie, timeout time.Duration) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Failed to call [%s], status code: %d", url, res.StatusCode)
	}

	return ioutil.ReadAll(res.Body)
}

func HTTPPost(url string, bodyData []byte) ([]byte, error) {
	b := bytes.NewReader(bodyData)
	res, err := http.Post(url, "plain/text", b)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Failed to call [%s], status code: %d", url, res.StatusCode)
	}

	return ioutil.ReadAll(res.Body)
}

func HTTPPostV1(url string, bodyData []byte, headers map[string]string, cookies []*http.Cookie, timeout time.Duration) ([]byte, error) {
	b := bytes.NewReader(bodyData)
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}

	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Failed to call [%s], status code: %d", url, res.StatusCode)
	}

	return ioutil.ReadAll(res.Body)
}

func HTTPDelete(url string, bodyData []byte) ([]byte, error) {
	return CustomRequest("DELETE", url, bodyData)
}

func HTTPPut(url string, bodyData []byte) ([]byte, error) {
	return CustomRequest("PUT", url, bodyData)
}

func HTTPParseURL(srcURL string, params map[string][]string) (string, error) {
	u, err := url.Parse(srcURL)
	if err != nil {
		return "", fmt.Errorf("Failed to parse url: %s", err.Error())
	}
	p := url.Values{}
	for k, v := range params {
		for _, param := range v {
			if len(p.Get(k)) > 0 {
				p.Add(k, param)
			} else {
				p.Set(k, param)
			}
		}
	}
	u.RawQuery = p.Encode()
	return u.String(), nil
}

// FormPost do a HTTP POST with x-www-form-urlencoded body
func FormPost(url string, bodyData []byte) ([]byte, error) {
	b := bytes.NewReader(bodyData)
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Failed to call [%s], status code: %d", url, res.StatusCode)
	}

	return ioutil.ReadAll(res.Body)
}
