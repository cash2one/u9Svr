package tool

import (
	"encoding/base64"
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
