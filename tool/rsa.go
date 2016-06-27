package tool

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
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

func ParsePkCS8PrivateKeyWithFile(file string) (key *rsa.PrivateKey, err error) {
	var keyBytes []byte
	keyBytes, err = DecodePemFile(file)
	if err != nil {
		return
	}
	ret, err := x509.ParsePKCS8PrivateKey(keyBytes)
	return ret.(*rsa.PrivateKey), err
}

func ParsePkCS8PrivateKeyWithStr(privateKey string) (key *rsa.PrivateKey, err error) {
	var keyBytes []byte
	keyBytes, err = base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return
	}
	ret, err := x509.ParsePKCS8PrivateKey(keyBytes)
	return ret.(*rsa.PrivateKey), err
}

func RsaPKCS1V15Sign(key *rsa.PrivateKey, hashType crypto.Hash, context string) (string, error) {
	hashType = crypto.SHA1

	if !hashType.Available() {
		return "", x509.ErrUnsupportedAlgorithm
	}

	h := hashType.New()

	//h := sha1.New()
	h.Write([]byte(context))
	digest := h.Sum(nil)

	bytes, err := rsa.SignPKCS1v15(nil, key, hashType, digest)
	ret := base64.StdEncoding.EncodeToString(bytes)

	return ret, err
}

func RsaVerifyPKCS1v15(key *rsa.PublicKey, hashType crypto.Hash, context, sign string) error {
	if !hashType.Available() {
		return x509.ErrUnsupportedAlgorithm
	}

	h := hashType.New()
	//h := sha1.New()
	h.Write([]byte(context))
	digest := h.Sum(nil)

	ds, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(key, hashType, digest, ds)
}
