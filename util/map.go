package util

func Map(x, inmin, inmax, outmin, outmax int64) int64 {
	return (x-inmin)*(outmax-outmin)/(inmax-inmin) + outmin
}
