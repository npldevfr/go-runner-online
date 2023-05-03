/*
//  Implementation of a few utilitary functions:
//    - getSeconds
*/

package main

import "strconv"

// GetSeconds splits a time in milliseconds into seconds and milliseconds
func GetSeconds(d int64) (s, ms int64) {
	s = d / 1000
	ms = d - s*1000
	return s, ms
}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
