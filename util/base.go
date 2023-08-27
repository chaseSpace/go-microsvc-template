package util

import (
	"golang.org/x/exp/rand"
	"time"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}
