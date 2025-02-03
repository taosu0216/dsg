package tools

import (
	"dsg/p"
	"encoding/json"
	"fmt"
	"os"
)

type AddDictReq struct {
	Filepath string `json:"filepath"`
}

func (t *Tool) CreateDict(info string) {
	fmt.Println("create dict info ", info)
	obj := &AddDictReq{}
	err := json.Unmarshal([]byte(info), obj)
	if err != nil {
		p.Red(fmt.Sprintf("json unmarshal err is %v", err))
		return
	}
	filePath := obj.Filepath

	// 调用 isExist 函数检查路径是否存在
	exists, _ := isExist(filePath)
	if exists {
		p.Red("dict is exist")
		return
	}

	err = os.MkdirAll(filePath, 0755)
	if err != nil {
		p.Red(fmt.Sprintf("create dict err is %v", err))
		return
	}

	p.Red("create dict successfully")
}
