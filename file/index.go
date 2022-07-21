package l_file

import (
	"bufio"
	"fmt"
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

func WriteTxt(filename string, content string) (err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)

	write.WriteString(content)

	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	return
}
