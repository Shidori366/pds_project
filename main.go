package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

func getTemplateAndSegmentCount(n int) ([]int, int) {
	m := int(math.Sqrt(float64(n))) + 1

	template := make([]int, m)
	p := 2

	for i := range template {
		if i < 2 {
			continue
		}
		template[i] = 1
	}

	for i := 2; i < m; i++ {
		multiplier := 2
		if template[i] == 1 {
			p = i
			for p*multiplier < m {
				template[multiplier*p] = 0
				multiplier += 1
			}
		}
	}

	segmentCount := n / m

	if n%m != 0 {
		segmentCount = n/m + 1
	}

	return template, segmentCount - 1
}

func processSegment(segmentStart int, template []int, length int) []int {
	segment := make([]int, length)

	for i := range segment {
		segment[i] = 1
	}

	for i, isPrime := range template {
		if isPrime == 0 {
			continue
		}

		multiplier := (segmentStart + i - 1) / i

		for i*multiplier-segmentStart < length {
			segment[i*multiplier-segmentStart] = 0
			multiplier++
		}
	}

	return segment
}

func process(n int) {
	template, segmentCount := getTemplateAndSegmentCount(n)
	l := len(template)
	nextSegmentStart := l

	for i := range segmentCount {
		nextSegmentStart = l * (i + 1)
		if i+1 == segmentCount && n-nextSegmentStart < l {
			l = n - nextSegmentStart + 1
		}
		processSegment(nextSegmentStart, template, l)
	}

}

func processParallel(n int) {
	template, segmentCount := getTemplateAndSegmentCount(n)
	l := len(template)
	nextSegmentStart := l

	wg := sync.WaitGroup{}

	for i := range segmentCount {
		nextSegmentStart = l * (i + 1)
		if i+1 == segmentCount && n-nextSegmentStart < l {
			l = n - nextSegmentStart + 1
		}

		wg.Add(1)
		go func(segmentStart int, l int) {
			defer wg.Done()
			processSegment(segmentStart, template, l)
		}(nextSegmentStart, l)
	}

	wg.Wait()
}

func main() {
	n := 0

	fmt.Print("Enter number: ")
	fmt.Scan(&n)

	var min int64 = 0
	var max int64 = 0
	var averageTime int64 = 0

	for i := range 30 {
		start := time.Now().UnixMilli()
		processParallel(n)
		end := time.Now().UnixMilli()

		currTime := end - start

		if i == 0 {
			min = currTime
			max = currTime
		}

		if min > currTime {
			min = currTime
		}

		if max < currTime {
			max = currTime
		}

		averageTime += end - start
	}

	averageTime /= 30
	fmt.Println("min: ", min, "ms")
	fmt.Println("max: ", max, "ms")
	fmt.Println("Average: ", averageTime, "ms")
}
