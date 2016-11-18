package math

import (
	"math"
)

//Inverse of normal distribution, convert from OC version
//See: https://stackoverflow.com/questions/1687497/objective-c-implementation-of-the-wilson-score-interval
func PNormalDist(qn float64) float64 {
	b := []float64{1.570796288, 0.03706987906, -0.8364353589e-3, -0.2250947176e-3, 0.6841218299e-5, 0.5824238515e-5, -0.104527497e-5, 0.8360937017e-7, -0.3231081277e-8, 0.3657763036e-10, 0.6936233982e-12}
	if qn < 0.0 || qn > 1.0 {
		return 0.0
	}

	if qn == 0.5 {
		return 0.0
	}
	w1 := qn
	if qn > 0.5 {
		w1 = 1.0 - w1
	}
	w3 := -math.Log(4.0 * w1 * (1.0 - w1))
	w1 = b[0]
	for i, v := range b[1:] {
		w1 += v * math.Pow(w3, float64(i+1))
	}

	if qn > 0.5 {
		return math.Sqrt(w1 * w3)
	}
	return -math.Sqrt(w1 * w3)
}

//Low bound of Wilson score confidence interval for a Bernoulli parameter
//See: http://www.evanmiller.org/how-not-to-sort-by-average-rating.html
func LBWilsonScoreWithBernoulliParam(positive int64, all int64, confidence float64) float64 {
	if all <= 0 {
		return 0
	}
	fPos := float64(positive)
	fAll := float64(all)
	z := PNormalDist(1.0 - (confidence)/2.0)
	phat := 1.0 * fPos / fAll
	return (phat + z*z/(2.0*fAll) - z*math.Sqrt((phat*(1-phat)+z*z/(4*fAll))/fAll)) / (1 + z*z/fAll)
}
