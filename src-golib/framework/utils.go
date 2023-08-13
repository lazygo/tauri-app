package framework

import "strings"

func cleanPath(path string) string {
	path = "/" + strings.Trim(path, "\\/ ")
	return strings.ToLower(path)
}
