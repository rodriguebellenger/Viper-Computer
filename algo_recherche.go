package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	for i := 1; i < 9; i++ {
		var list []int = []int{}
		var limit float64 = math.Pow(10, float64(i))
		for j := 0; j < int(limit); j++ {
			list = append(list, j)
		}
		var startTime time.Time = time.Now()
		search(list, -1)
		var elapsed time.Duration = time.Since(startTime)
		fmt.Printf("%d : %s\n", int(limit), elapsed)

		startTime = time.Now()
		binarySearch(list, -1)
		elapsed = time.Since(startTime)
		fmt.Printf("%d : %s\n", int(limit), elapsed)
	}
}

func search(list []int, x int) bool {
	if len(list) == 0 {
		return false
	}
	for _, element := range list {
		if element == x {
			return true
		}
	}
	return false
}

func binarySearch(list []int, x int) bool {
	var start int = 0
	var end int = len(list) - 1
	var middle int = (len(list) - 1) / 2
	for start <= end {
		if list[middle] < x {
			start = middle + 1
		} else if list[middle] > x {
			end = middle - 1
		} else {
			return true
		}
		middle = (start + end) / 2
	}
	return false
}
