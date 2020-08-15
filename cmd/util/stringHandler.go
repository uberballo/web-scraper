package util

import (
	"regexp"
	"strings"
)

//AppendSuffix appends suffix
func AppendSuffix(list []string, suffix string) []string {
	var res []string
	for _, n := range list {
		res = append(res, n+suffix)
	}
	return res
}

//PrependPrefix prepends prefix
func PrependPrefix(list []string, prefix string) []string {
	var res []string
	for _, n := range list {
		res = append(res, prefix+n)
	}
	return res

}

func removePartOfString(url, toRemove string) string {
	res := strings.Replace(url, toRemove, "", -1)
	return res
}

//GetLastPart Gets the last part of a url.
// Eg. example/hello/part -> part
func GetLastPart(url string) string {
	re := regexp.MustCompile(`([^\/]+$)`)
	trimmedURL := removePartOfString(url, "/tilinpaatos")
	symbol := re.Find([]byte(trimmedURL))
	return string(symbol)
}
