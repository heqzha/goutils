package test

import (
	"fmt"
	"github.com/heqzha/goutils/date"
	"github.com/heqzha/goutils/math"
	"testing"
)

func TestPNormalDist(t *testing.T) {
	fmt.Printf("PNormalDist: %g\n", math.PNormalDist(0.95))
}

func TestLBWilsonScoreWithBernoulliParam(t *testing.T) {
	fmt.Printf("LBWilsonScoreWithBernoulliParam: %g\n", math.LBWilsonScoreWithBernoulliParam(500, 1000, 0.95))
}

func TestRedditHotRankScore(t *testing.T) {
	ups, downs, delta := int64(100), int64(50), int64(500)
	s1 := math.RedditHotRankScore(ups, downs, date.DateNowSecond())
	t.Log(s1)

	ups = ups + delta
	s2 := math.RedditHotRankScore(ups, downs, date.DateNowSecond())
	t.Log(s2)
}
