package wueste

import "strings"

func ArrayLessString(a, b []string) bool {
	if len(a) < len(b) {
		return true
	}
	for i := 0; i < len(a); i++ {
		c := strings.Compare(a[i], b[i])
		if c < 0 {
			return true
		}
		if c > 0 {
			return false
		}
	}
	return false
}

func ArrayLessInteger[T uint | int | uint64 | uint32 | uint16 | uint8 | int8 | int16 | int32 | int64](a, b []T) bool {
	if len(a) < len(b) {
		return true
	}
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return true
		}
	}
	return false
}
func ArrayLessNumber[T float32 | float64](a, b []T) bool {
	if len(a) < len(b) {
		return true
	}
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return true
		}
	}
	return false
}
func ArrayLessBoolean(a, b []bool) bool {
	if len(a) < len(b) {
		return true
	}
	for i := 0; i < len(a); i++ {
		if a[i] == false && b[i] == true {
			return true
		}
		if a[i] != b[i] {
			return false
		}
	}
	return false
}
