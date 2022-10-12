package utils

import (
	"strings"
)

// parse service path
func ParseServicePath(path string) (string, string) {
	arr := strings.Split(path, ".")
	return strings.Join(arr[:len(arr)-1], "."), arr[len(arr)-1]
}
