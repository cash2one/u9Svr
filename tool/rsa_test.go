package tool

import (
	"crypto"
	"fmt"
	"testing"
)

func TestRsaVerify(t *testing.T) {
	var hashType crypto.Hash
	hashType = crypto.MD5SHA1
	publicKey := `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCeZ6HFJXiXORcv5ljI27J8ZUb/YIXDzRIpVN53NOgZ0NZ4OplXPumZBxR/gksskd79sPMcy9Rvpz8ZiPUKTTUuTmUMjtL9f/E1XafVcjvUUrUILv+aJb65OiR9YHqbGSqj8B9qR5pmtyP8TAuBA2CRooBF01WrYRHXxYv328aDWwIDAQAB`
	text := "notifyId=20160204152757964736&partnerOrder=20160204152757165&productName=test&productDesc=test&price=5&count=1&attach=test&sign=8f00a109716e819bfe0afb695c1addf1"
	sign := "BV52daLE+DCRCByIUvu9SpxXC9ov/ftHp7aWiOwu/lBA2FapBl2akxT1MYMinxEZTf4VrJbhnHC8/pPHo5nY4EGykHOmJk6AXm8GwgYlk7AK5O9wUSqA+61UD0OlefNyuCuVuVQabDEu0RS6Q99D2mN99M5ALOJODDWC4GOShNE="
	keyInterface, err := ParsePKIXPublicKeyWithStr(publicKey)
	if err != nil {
		fmt.Println("ParsePKIXPublicKeyWithStr:", err)
		return
	}
	err = RsaVerifyPKCS1v15(keyInterface, hashType, text, sign)
	fmt.Println("RsaVerifyPKCS1v15WithStr:", err)
}
