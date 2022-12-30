package utils

import (
	"fmt"
	"strconv"

	"github.com/Knetic/govaluate"
)

func GetSum(nums ...float64) (sum float64) {
	for _, v := range nums {
		sum += v
	}
	return
}

func GetMean(nums ...float64) float64 {
	sum := GetSum(nums...)
	n := len(nums)
	return sum / float64(n)
}

// Variance: s^2 = ((x1 - m)^2 + (x2 - m)^2 + ... + (xn - m)^2) / n
func GetVariance(nums ...float64) float64 {
	e := fmt.Sprintf("(%s) / n", func(nums ...float64) (s string) {
		for index := range nums {
			if index == 0 {
				s = fmt.Sprintf("(x%s - m)**2", strconv.Itoa(index+1))
				continue
			}
			s = fmt.Sprintf("%s + (x%s - m)**2", s, strconv.Itoa(index+1))
		}
		return
	}(nums...))
	expr, _ := govaluate.NewEvaluableExpression(e)
	parameters := make(map[string]interface{})

	for index, value := range nums {
		p := fmt.Sprintf("x%s", strconv.Itoa(index+1))
		parameters[p] = value
	}
	parameters["m"] = GetMean(nums...)
	parameters["n"] = len(nums)

	result, _ := expr.Evaluate(parameters)
	return result.(float64)
}
