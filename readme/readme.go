package readme

import (
	"dsg/p"
	"fmt"
)

func Hi() {
	// ASCII 艺术打印项目名称
	p.Hi(" ____     _____     ____ ")
	p.Hi(`|  _ \   |  ___|   / ___|`)
	p.Hi("| | | |  | |___   | |  _ ")
	p.Hi("| |_| |   ____ |  | |_| |")
	p.Hi(`|____/   |_____|   \____|`)
	fmt.Println()

	// 项目简介
	p.Hi("这是一个基于 Golang 开发的命令行工具，旨在提供个性化的cli级别的 Cursor 的诸如ai生成，编辑文件，创建目录，聊天等功能。")
	p.Hi("核心功能包括：")
	p.Hi("- 文件创建与写入：通过 /addFile 命令，可以快速创建文件并写入指定内容。")
	p.Hi("- 目录创建：通过 /addDict 命令，可以一键创建指定目录，简化文件管理。")
	p.Hi("- 文件编辑：通过 /editFile 命令，可以灵活编辑文件内容，支持内容替换和删除操作。")
	p.Hi("- 命令行级别的 Cursor 或 Windsurf 功能：支持用户通过命令行交互式管理文件和数据，提供流畅的操作体验。")
	p.Hi("- AI 辅助编辑：通过内置的 AI 功能，可以智能地帮助用户编辑文件内容，提高工作效率。")
	p.Hi("使用方法：")
	p.Hi("1. 添加文件到上下文(可以输入多个)：/addFile <filepath>")
	p.Hi("2. 添加目录到上下文(可以输入多个)：/addDict <dictpath>")
	p.Hi("3. 准备让ai帮忙 创建/编辑文件，创建目录，使用/tool 回车之后直接说需求即可(普通聊天是流式响应，目前工具调用相关无法流式响应)")
	p.Hi("4. 退出程序：/quit")
}
