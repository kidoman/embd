// Package util contains utility functions.
package util

// Map re-maps a number from one range to another.
//
// Example:
//
//	val := Map(angle, 0, 180, 1000, 2000)
//
func Map(x, inmin, inmax, outmin, outmax int64) int64 {
	return (x-inmin)*(outmax-outmin)/(inmax-inmin) + outmin
}
