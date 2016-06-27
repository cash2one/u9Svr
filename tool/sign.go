package tool

import (
	//"fmt"
	"os/exec"
	"strings"
)

func IapppaySign(transdata, key string) (result string, err error) {
	var args []string
	args = make([]string, 4)
	args[0] = "-jar"
	args[1] = "../tool/jar/iapppay_sign.jar"
	args[2] = transdata
	args[3] = key
	var buf []byte
	cmd := exec.Command("java", args...)
	buf, err = cmd.Output()
	result = string(buf)
	return
}

func IapppayVerify(transdata, sign, key string) (result string, err error) {
	var args []string
	args = make([]string, 5)
	args[0] = "-jar"
	args[1] = "../tool/jar/iapppay_verify.jar"
	args[2] = transdata
	args[3] = sign
	args[4] = key
	var buf []byte
	cmd := exec.Command("java", args...)
	buf, err = cmd.Output()
	result = string(buf)
	return
}

func YYHSign(transdata, sign, key string) (result string, err error) {
	var args []string

	args = make([]string, 5)
	args[0] = "-jar"
	args[1] = "../tool/jar/yyh_sign622.jar"
	args[2] = transdata
	args[3] = sign
	args[4] = key

	var buf []byte
	cmd := exec.Command("java", args...)
	buf, err = cmd.Output()
	result = string(buf)
	return
}

func MMYSign(sign, appKey, orderId string) (result string, err error) {
	var args []string
	args = make([]string, 5)
	args[0] = "-jar"
	args[1] = "../tool/jar/mmy_sign.jar"
	args[2] = sign
	args[3] = appKey
	args[4] = orderId
	var buf []byte
	cmd := exec.Command("java", args...)
	buf, err = cmd.Output()
	result = string(buf)
	// fmt.Println(result)
	return
}

func TTSign(signContent string) (result string, err error) {
	var args []string
	args = make([]string, 5)
	args[0] = "-jar"
	args[1] = "../tool/jar/tt_sign.jar"
	args[2] = signContent
	var buf []byte
	cmd := exec.Command("java", args...)
	buf, err = cmd.Output()
	result = string(buf)
	result = strings.TrimSpace(result)
	// fmt.Println(result)
	return
}
