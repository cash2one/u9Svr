package tool

import (
	"os/exec"
)

func UnTar(file string, path string) (err error) {
	var args []string
	args = make([]string, 4)
	args[0] = "-xvf"
	args[1] = file
	args[2] = "-C"
	args[3] = path

	cmd := exec.Command("tar", args...)
	//var buf []byte
	//buf, err = cmd.Output()
	_, err = cmd.Output()
	//if err != nil {
	//fmt.Fprintf(os.Stderr, "The command failed to perform: %s", err)
	//}
	//result = string(buf)
	//fmt.Fprintf(os.Stdout, "Result: %s", buf)
	return
}

func CopyDir(srcDir string, destDir string) (err error) {
	var args []string
	args = make([]string, 3)
	args[0] = "-Rvf"
	args[1] = srcDir
	args[2] = destDir

	cmd := exec.Command("cp", args...)
	_, err = cmd.Output()
	return
}

func CreateDir(dirName string) (err error) {
	var args []string
	args = make([]string, 2)
	args[0] = "-p"
	args[1] = dirName

	cmd := exec.Command("mkdir", args...)
	_, err = cmd.Output()
	return err
}
