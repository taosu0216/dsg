package tools

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Tool struct {
}

func NewTool() *Tool {
	return &Tool{}
}

type ToolAdd interface {
	Add(string) (string, error)
}

func isExist(path string) (bool, string) {
	// 使用 os.Stat 检查路径是否存在
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 如果路径不存在，返回 false 和空字符串
		return false, ""
	} else if err != nil {
		// 如果发生其他错误，返回 false 和错误信息
		return false, err.Error()
	}

	if fileInfo.IsDir() {
		return true, Dict
	} else {
		return true, File
	}
}

func (t *Tool) IsExistAndGetContent(url string) (bool, string, error) {
	if exist, t := isExist(url); exist {
		if t == File {
			fmt.Println(url, "   here  ")
			contentBytes, err := ioutil.ReadFile(url)
			if err != nil {
				log.Println("read file fail ", err)
				return true, "", errors.New("internal err")
			}
			content := string(contentBytes)
			//fmt.Println("fule  ",content)
			return true, content, nil
		} else {
			return true, "", errors.New("not a file")
		}
	} else {
		return false, "", nil
	}
}
