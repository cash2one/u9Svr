package common

import (
	"u9/tool"
)

const (
	passwordSignKey = "shine"
)

func GetAuthKey(name string) (ret string) {
	return tool.Md5([]byte(name + passwordSignKey))
}
