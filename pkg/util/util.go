package util

import (
	"strconv"
	"time"
)

func MustAtoi(input string) int {
	ans, _ := strconv.Atoi(input)
	return ans
}

func StringsToInts(input []string) ([]int, error) {
	var ints []int
	for _, val := range input {
		tmp, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		ints = append(ints, tmp)
	}
	return ints, nil
}

func StringsToTimes(input []string) ([]time.Time, error) {
	var times []time.Time
	layout := "2006-01-02T15:04"

	for _, val := range input {
		parsedTime, err := time.Parse(layout, val)
		if err != nil {
			return nil, err
		}
		times = append(times, parsedTime)
	}
	return times, nil
}
