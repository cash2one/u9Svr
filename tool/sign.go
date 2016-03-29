package tool

import (
	"os/exec"
)

func YYHSign(transdata, sign, key string) (result string, err error) {
	var args []string
	args = make([]string, 5)
	args[0] = "-jar"
	args[1] = "../tool/jar/yyh_sign.jar"
	args[2] = transdata
	args[3] = sign
	args[4] = key
	var buf []byte
	cmd := exec.Command("java", args...)
	buf, err = cmd.Output()
	result = string(buf)
	// fmt.Println(result)
	return
}
