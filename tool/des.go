package tool

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"errors"
	"fmt"
)

func test() {
	context := `{"payAmount":"1000","uid":"20150306150802DH2jqPw0p6","notifyTime":1448446419,"cpInfo":"anzhi_20151125180905927","memo":null,"orderAmount":"1000","pt":"ZFB","orderAccount":"15385218158","code":1,"orderTime":"2014-12-03 00:03:07","msg":"","orderId":"1511251809950000002"}`
	key := "ug2KMdLi2JSr4naOE48XmL3h"
	//dest := `bdptFOVVuyMvvaK7ozVbq3x+zM/GVkerIYh9mgjjAWM36t26EpU3yknsaoxI6V9+XSfr4eZsJejfOqgWD+AzZ/t1UVrHC7c6KmsCj2Baw5/nz5sB2OE5tpgEx35TxzYLft29toIt/wFqrHiId0Z9uEZZicKjWOG686P5vwB3+fq9t+nbw/fcWJ81gXe2uvfyBJNOt1p/Ve/pO1S4Fk3RhTqNayxKBp3e1Gf2HLCoMla5f9YtzUvTXR8cHksbHzyLtlnwzsmCblDDqe3p/KmccCLZuf+9+1JiYZlbgUmsjoCdCHBoRiZPMenx67XIvc931+qQb0rSgyFwI9EM6e4mgTqD3x6FGjaOxwl/bwNQdek=`
	dest, _ := JavaDesEncyrpt(key, context)
	fmt.Println(dest)
	src, _ := JavaDesDecyrpt(key, dest)
	fmt.Println(src)
}

func JavaDesEncyrpt(key, context string) (ret string, err error) {
	return DesEncrypt("3des", "pkcs5Padding", "ecb", key, context)
}

func JavaDesDecyrpt(key, context string) (ret string, err error) {
	return DesDecrypt("3des", "pkcs5Padding", "ecb", key, context)
}

func DesEncrypt(algorithm, paddingMode, encryptMode, key, context string) (ret string, err error) {
	data := []byte(context)
	var block cipher.Block
	switch algorithm {
	case "des":
		block, err = des.NewCipher([]byte(key))
	case "3des":
		block, err = des.NewTripleDESCipher([]byte(key))
	}
	if err != nil {
		return
	}

	bs := block.BlockSize()
	switch paddingMode {
	case "pkcs5Padding":
		data = PKCS5Padding(data, bs)
	case "zeroPadding":
		data = ZeroPadding(data, bs)
	}
	if len(data)%bs != 0 {
		err = errors.New("Need a multiple of the blocksize")
		return
	}

	out := make([]byte, len(data))
	switch encryptMode {
	case "ecb":
		dst := out
		for len(data) > 0 {
			block.Encrypt(dst, data[:bs])
			data = data[bs:]
			dst = dst[bs:]
		}
	case "cbc":
		blockMode := cipher.NewCBCEncrypter(block, []byte(key))
		blockMode.CryptBlocks(out, data)
	}
	return base64.StdEncoding.EncodeToString(out), nil
}

func DesDecrypt(algorithm, paddingMode, encryptMode, key, context string) (ret string, err error) {
	data, err := base64.StdEncoding.DecodeString(context)
	if err != nil {
		return
	}

	keyByte := []byte(key)
	var block cipher.Block
	switch algorithm {
	case "des":
		block, err = des.NewCipher(keyByte)
	case "3des":
		block, err = des.NewTripleDESCipher(keyByte)
		keyByte = keyByte[:8]
	}
	if err != nil {
		return
	}

	bs := block.BlockSize()
	if len(data)%bs != 0 {
		err = errors.New("crypto/cipher: input not full blocks")
		return
	}

	out := make([]byte, len(data))
	switch encryptMode {
	case "ecb":
		dst := out
		for len(data) > 0 {
			block.Decrypt(dst, data[:bs])
			data = data[bs:]
			dst = dst[bs:]
		}
	case "cbc":
		blockMode := cipher.NewCBCDecrypter(block, keyByte)
		blockMode.CryptBlocks(out, data)
	}

	switch paddingMode {
	case "pkcs5Padding":
		out = PKCS5UnPadding(out)
	case "zeroPadding":
		out = ZeroUnPadding(out)
	}
	return string(out), nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
