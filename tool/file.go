package tool

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func GetDirList(path string, extName string) (map[string]string, error) {
	ret := make(map[string]string, 0)
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	//PthSep := string(os.PathSeparator)
	extName = strings.ToUpper(extName)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), extName) {
			ret[fi.Name()] = path + "/" + fi.Name()
		}
	}
	return ret, nil
}

func CopyFile(srcName, dstName string) (int64, error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

