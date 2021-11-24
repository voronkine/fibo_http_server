package main

import (
	"testing"
)

func TestFibonachi(t *testing.T) {

	redis := newRedisClient()
	service := Service{
		redis: redis,
	}

	x, y := 2, 5

	result := service.fibonachi(x, y)

	want := []int{1, 1, 2, 3}

	if len(result) != len(want) {
		t.Errorf("fibonachi() working wrong (len)")
		return
	}

	for i := 0; i < len(result); i++ {
		if result[i] != want[i] {
			t.Errorf("fibonachi() working wrong (val)")
			return
		}
	}

}
