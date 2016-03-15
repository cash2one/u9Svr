package common

import (
	"net/url"
	"strings"
)

func Rawurlencode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}
