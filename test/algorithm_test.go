package test

import(
	"fmt"
	"testing"
	"github.com/heqzha/goutils/math"
)

func TestPNormalDist(t *testing.T){
	fmt.Printf("PNormalDist: %g\n", math.PNormalDist(0.95))
}

func TestLBWilsonScoreWithBernoulliParam(t *testing.T){
	fmt.Printf("LBWilsonScoreWithBernoulliParam: %g\n", math.LBWilsonScoreWithBernoulliParam(500, 1000, 0.95))
}
