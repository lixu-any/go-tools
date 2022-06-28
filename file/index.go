package l_file

import (
	"io/ioutil"
	"os"
)

//读取文本文件
func ReadFile(pth string) ([]byte, error) {

	filePtr, err := os.Open(pth)
	if err != nil {
		return nil, err
	}
	defer filePtr.Close()

	by, err := ioutil.ReadAll(filePtr)
	if err != nil {
		return nil, err
	}

	return by, nil

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//basePath是固定目录路径
func CreateDateDir(basePath string) (err error) {

	err = os.MkdirAll(basePath, 0777)
	if err != nil {
		return
	}

	return

}
