package tool

import (
	"encoding/base64"
	"strconv"
	"strings"
)

func DecodeStdEncodeing(context string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(context)
	if err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
	return "", nil
}

func Unicode2utf8(unicodeStr string) (utf8Str string) {
	//unicodeStr = `\u5546\u54c1\u63cf\u8ff0\`
	strlen := len(unicodeStr)
	if strings.HasSuffix(unicodeStr, `\`) {
		unicodeStr = unicodeStr[0 : strlen-1]
	}
	utf8Str, _ = strconv.Unquote(`"` + unicodeStr + `"`)
	return
}
