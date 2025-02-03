package tools

import (
	"dsg/p"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type EditFileRequest struct {
	Filepath   string `json:"filepath"`
	OldContent string `json:"oldContent"`
	NewContent string `json:"newContent"`
}

func (t *Tool) EditFile(info string) {
	fmt.Println("edit file info ", info)
	obj := &EditFileRequest{}
	err := json.Unmarshal([]byte(info), obj)
	if err != nil {
		p.Red(fmt.Sprintf("json unmarshal err is %v", err))
		return
	}
	path := obj.Filepath
	original := obj.OldContent
	newContent0 := obj.NewContent

	// 检查文件是否存在
	if exists, _ := isExist(path); !exists {
		p.Red("file isn't exist")
		return
	}

	// 读取文件内容
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		p.Red(fmt.Sprintf("read file err is %v", err))
		return
	}
	content := string(contentBytes)

	// 检查原始片段是否存在
	if !strings.Contains(content, original) {
		p.Red("original content isn't exist")
		return
	}

	// 替换第一个匹配项
	newContent := strings.Replace(content, original, newContent0, 1)

	// 写入文件
	if err := ioutil.WriteFile(path, []byte(newContent), 0644); err != nil {
		p.Red(fmt.Sprintf("replace file err is %v", err))
		return
	}

	p.Red("edit file successfully")
}
