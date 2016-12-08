package test

import (
	"fmt"
	"github.com/heqzha/goutils/unique"
	"testing"
)

func TestRandStringBytesMaskImprSrc(t *testing.T) {
	for i := 0; i < 10; i = i + 1 {
		c := unique.RandStringBytesMaskImprSrc(8)
		t.Log(fmt.Sprintf("Code %d: %s\n", i, c))
	}
}
