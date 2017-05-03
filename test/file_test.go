package test

import (
	"testing"

	"github.com/heqzha/goutils/file"
)

func TestGetSHA1(t *testing.T) {
	sha1, err := file.GetSHA1("./file_test.go")
	if err != nil {
		t.Error(err)
	}
	t.Log(len(sha1), sha1)
}
