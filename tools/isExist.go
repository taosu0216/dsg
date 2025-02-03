package tools

import (
	"dsg/p"
	"encoding/json"
	"fmt"
)

type IsExistReq struct {
	Filepath string `json:"filepath"`
}

func (t *Tool) IsExist(info string) bool {
	obj := &IsExistReq{}
	err := json.Unmarshal([]byte(info), obj)
	if err != nil {
		p.Red(fmt.Sprintf("json unmarshal err is %v", err))
		return true
	}
	filePath := obj.Filepath
	// 调用 isExist 函数检查路径是否存在
	exists, _ := isExist(filePath)
	if exists {
		p.Red(fmt.Sprintf("[%s] is exist", filePath))
		return true
	} else {
		p.Red(fmt.Sprintf("[%s] isn't exist", filePath))
		return false
	}
}
