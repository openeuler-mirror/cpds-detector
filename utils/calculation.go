package utils

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
