package tools

import (
	"dsg/p"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type AddFileReq struct {
	Filepath string `json:"filepath"`
	Content  string `json:"content"`
}

// createFile 创建文件并写入内容
func (t *Tool) CreateFile(info string) {
	fmt.Println("create file info ", info)
	obj := &AddFileReq{}
	err := json.Unmarshal([]byte(info), obj)
	if err != nil {
		p.Red(fmt.Sprintf("json unmarshal err is %v", err))
		return
	}
	filePath := obj.Filepath
	content := obj.Content

	// 调用 isExist 函数检查路径是否存在
	exists, _ := isExist(filePath)
	if exists {
		p.Red("file is exist")
		return
	}

	// 如果文件不存在，创建文件并写入内容
	err = ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		p.Red(fmt.Sprintf("write to file err is %v", err))
		return
	}

	p.Red("add field successfully")
}
