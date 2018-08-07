package dpos

import (
	"fmt"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	seed := uint64(time.Now().Unix())
	if seed <= 0 {
		t.Errorf("seed less than zero.")
	}

	rand := NewRandom(seed)

	for i := 0; i < 5; i++ {
		fmt.Println(rand.GenRandom())
	}
}
