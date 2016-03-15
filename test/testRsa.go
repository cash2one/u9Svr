package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func DecodePemFile(file string) (keyBytes []byte, err error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0400)
	if err != nil {
		return
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	p, _ := pem.Decode(buf)
	if p == nil {
		return nil, errors.New("no pem block found")
	}
	keyBytes = p.Bytes
	return
}

func ParsePKIXPublicKeyWithFile(file string) (key *rsa.PublicKey, err error) {
	var keyBytes []byte
	keyBytes, err = DecodePemFile(file)
	if err != nil {
		return
	}

	var keyInterface interface{}
	keyInterface, err = x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return
	}
	key = keyInterface.(*rsa.PublicKey)
	return
}

func ParsePKIXPublicKeyWithStr(publicKey string) (key *rsa.PublicKey, err error) {
	var keyBytes []byte
	keyBytes, err = base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return
	}

	var keyInterface interface{}
	keyInterface, err = x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return
	}
	key = keyInterface.(*rsa.PublicKey)
	return
}

func ParsePkCS1PrivateKeyWithFile(file string) (key *rsa.PrivateKey, err error) {
	var keyBytes []byte
	keyBytes, err = DecodePemFile(file)
	if err != nil {
		return
	}
	return x509.ParsePKCS1PrivateKey(keyBytes)
}

func ParsePkCS1PrivateKeyWithStr(privateKey string) (key *rsa.PrivateKey, err error) {
	var keyBytes []byte
	keyBytes, err = base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return
	}
	return x509.ParsePKCS1PrivateKey(keyBytes)
}

func RsaVerifyPKCS1v15(key *rsa.PublicKey, context, sign string) error {
	h := sha1.New()
	h.Write([]byte(context))
	digest := h.Sum(nil)

	ds, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(key, crypto.SHA1, digest, ds)
}

func main() {

	publicKey := `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCeZ6HFJXiXORcv5ljI27J8ZUb/YIXDzRIpVN53NOgZ0NZ4OplXPumZBxR/gksskd79sPMcy9Rvpz8ZiPUKTTUuTmUMjtL9f/E1XafVcjvUUrUILv+aJb65OiR9YHqbGSqj8B9qR5pmtyP8TAuBA2CRooBF01WrYRHXxYv328aDWwIDAQAB`
	//key, _ := base64.StdEncoding.DecodeString(publicKey)
	//re, err := x509.ParsePKIXPublicKey(key)
	//pub := re.(*rsa.PublicKey)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	text := "notifyId=20160204152757964736&partnerOrder=20160204152757165&productName=test&productDesc=test&price=5&count=1&attach=test&sign=8f00a109716e819bfe0afb695c1addf1"
	sign := "BV52daLE+DCRCByIUvu9SpxXC9ov/ftHp7aWiOwu/lBA2FapBl2akxT1MYMinxEZTf4VrJbhnHC8/pPHo5nY4EGykHOmJk6AXm8GwgYlk7AK5O9wUSqA+61UD0OlefNyuCuVuVQabDEu0RS6Q99D2mN99M5ALOJODDWC4GOShNE="
	//h := sha1.New()
	//h.Write(text)
	//digest := h.Sum(nil)

	//ds, _ := base64.StdEncoding.DecodeString(sig)
	//err = rsa.VerifyPKCS1v15(pub, crypto.SHA1, digest, ds)
	keyInterface, err := ParsePKIXPublicKeyWithStr(publicKey)
	if err != nil {
		fmt.Println("ParsePKIXPublicKeyWithStr:", err)
		return
	}
	err = RsaVerifyPKCS1v15(keyInterface, text, sign)
	fmt.Println("RsaVerifyPKCS1v15WithStr:", err)
}

/*package main

import (
	. "crypto"
	"encoding/base64"
	"fmt"
	"github.com/89hmdys/toast/crypto"
	"github.com/89hmdys/toast/rsa"
	//"testing"
)

func main() {

	plant := `test`
	//plant = "partnerOrder=20160204152757165&notifyId=20160204152757964736" +
	//	"productName=test&productDesc=test&price=5&count=1&attach=test"
	//plant:=`G|U09`

	key, err := rsa.LoadPKCS8KeyFromPEMFile(
		`d:/go/src/github.com/89hmdys/toast/crypto/rsa_public_key.pem`,
		`d:/go/src/github.com/89hmdys/toast/crypto/rsa_private_key.pem`)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	cipher, err := crypto.NewRSA(key)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	enT, err := cipher.Encrypt([]byte(plant))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Println(base64.StdEncoding.EncodeToString(enT))

	deT, err := cipher.Decrypt(enT)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Println(string(deT))
}*/
