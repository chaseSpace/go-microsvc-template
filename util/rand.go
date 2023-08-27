package util

import "golang.org/x/exp/rand"

func RandIntRange(left, right int, duplicate map[int]int) int {
	if left >= right {
		return left
	}
	delta := right - left
	if len(duplicate) == delta+1 { // full
		return 0
	}
	for {
		i := rand.Intn(delta+1) + left
		if duplicate == nil || duplicate[i] == 0 {
			duplicate[i] = 1
			return i
		}
	}
}
