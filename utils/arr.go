package utils

func ArrSum(arr []float64) float64 {
	var sum float64
	idx := 0
	for {
		if idx > len(arr)-1 {
			break
		}
		sum += arr[idx]
		idx++
	}
	return sum
}
