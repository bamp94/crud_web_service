package main

import (
	"fmt"
	"sort"
)

func getSumCombinations(combinations []int, sum int) [][]int {
	if len(combinations) < 3 {
		return nil
	}
	sort.Ints(combinations)
	var sumCombinations [][]int
	for i := len(combinations) - 1; i > 2; i-- {
		first := 0
		penult := i - 1
		last := i
		for first < penult {
			s := combinations[first] + combinations[penult] + combinations[last]
			if s == sum {
				sumCombinations = append(sumCombinations, []int{combinations[first], combinations[penult], combinations[last]})
			}

			if s < sum {
				first += 1
			} else {
				penult -= 1
			}
		}
	}
	return sumCombinations
}

func main() {
	fmt.Printf("Sum = %v, Combinations: %v\n", 213123123, getSumCombinations([]int{-1, 2, 22, 213123123, 0, 1}, 213123123))
	fmt.Printf("Sum = %v, Combinations: %v\n", 6, getSumCombinations([]int{-1, 2, 4, 3, 0, 1}, 6))
}
