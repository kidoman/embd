package util

import "testing"

func TestMap(t *testing.T) {
	var tests = []struct {
		x, inmin, inmax, outmin, outmax int64
		val                             int64
	}{
		{
			90, 0, 180, 1000, 2000,
			1500,
		},
		{
			10, 10, 15, 10, 20,
			10,
		},
		{
			15, 10, 15, 10, 20,
			20,
		},
	}
	for _, test := range tests {
		val := Map(test.x, test.inmin, test.inmax, test.outmin, test.outmax)
		if val != test.val {
			t.Errorf("Map of %v from (%v -> %v) to (%v -> %v): got %v, want %v", test.x, test.inmin, test.inmax, test.outmin, test.outmax, val, test.val)
		}
	}
}
