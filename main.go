package main

import (
	"bufio"
	"context"
	"dsg/p"
	"dsg/prompt"
	"dsg/readme"
	"dsg/tools"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/taosu0216/deepseek"
)

var (
	cli    *deepseek.Client
	tool   *tools.Tool
	dialog []deepseek.Message
)

func init() {
	dialog = make([]deepseek.Message, 0)
	tool = tools.NewTool()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	} else {
		// 从环境变量中获取assPropmt的值
		assPropmt := os.Getenv("PROMPT")
		if assPropmt != "" {
			prompt.Assistant = assPropmt
			p.Yellow("your prompt is: ")
			fmt.Println(assPropmt)
		}
		cli = deepseek.NewClient(os.Getenv("KEY")).WithBaseUrl(os.Getenv("URL"))
		fmt.Println(os.Getenv("KEY"), "    ", os.Getenv("URL"))
		fmt.Println()
		fmt.Println()
		fmt.Println("------------------------------------------------------------------------------------")
		fmt.Println()
		fmt.Println()
	}

	readme.Hi()
}

func main() {
	ctx := context.Background()

	addFileParams := deepseek.NewParameters().
		WithProperty("filepath", deepseek.ToolParamTypeStr, "文件路径，例如：/root/gopro/1.txt").
		WithProperty("content", deepseek.ToolParamTypeStr, "文件内容，例如: hello world").
		WithRequired("filepath", "content")

	addDictParams := deepseek.NewParameters().
		WithProperty("filepath", deepseek.ToolParamTypeStr, "路径，例如：/root/gopro/tmp").
		WithRequired("filepath")

	editFileParams := deepseek.NewParameters().
		WithProperty("filepath", deepseek.ToolParamTypeStr, "文件路径，例如：/root/gopro/1.txt").
		WithProperty("oldContent", deepseek.ToolParamTypeStr, "旧的需要替换的文件内容，例如: hello world").
		WithProperty("newContent", deepseek.ToolParamTypeStr, "新的文件内容，例如: 你好世界").
		WithRequired("filepath", "oldContent", "newContent")

	isExistParams := deepseek.NewParameters().
		WithProperty("filepath", deepseek.ToolParamTypeStr, "文件路径，例如：/root/gopro/1.txt").
		WithRequired("filepath")

	cliTools := []*deepseek.Tool{
		deepseek.NewTool("addFile", "创建文件并写入内容", addFileParams),
		deepseek.NewTool("addDict", "创建一个目录", addDictParams),
		deepseek.NewTool("editFile", "编辑一个文件，对指定旧的内容替换成新内容，即可以修改文件，也可以删除文件对应内容，比如把旧内容替换成空", editFileParams),
		deepseek.NewTool("isExist", "判断文件或目录是否存在", isExistParams),
	}

	dialog = append(dialog, []deepseek.Message{
		{
			Role:    deepseek.ChatMessageRoleSystem,
			Content: prompt.SystemCons,
		},
		{
			Role:    "assistant",
			Content: prompt.Assistant,
		},
	}...)
	baseReq := deepseek.ChatCompletionRequest{
		Model:                    os.Getenv("MODEL"),
		Messages:                 dialog,
		ResponseType:             deepseek.JSON_OBJECT,
		Tools:                    cliTools,
		ChatCompletionToolChoice: "auto",
	}
	quitCh := make(chan struct{})
	go func() {
		<-quitCh
		p.Green("chat stop")
		os.Exit(0)
	}()

	p.Red("--------------  strat chat --------------")
	lastLineEmpty := false
	toolFlag := false
	your()

	scanner := bufio.NewScanner(os.Stdin)
	var input strings.Builder
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "/addFile ") {
			remaining := strings.TrimPrefix(line, "/addFile ")
			realFiles := strings.Fields(remaining)
			addFileToCtx(realFiles)
			your()
			continue
		} else if strings.HasPrefix(line, "/addDict ") {
			remaining := strings.TrimPrefix(line, "/addDict ")
			realFiles := strings.Fields(remaining)
			addDictToCtx(realFiles)
			your()
			continue
		} else if strings.HasPrefix(line, "/tool") {
			toolFlag = true
			your()
			continue
		} else if strings.HasPrefix(line, "/quit") {
			quitCh <- struct{}{}
			continue
		}

		if line == "" {
			if lastLineEmpty {
				ai()
				allInfo := input.String()
				dialog = append(dialog, deepseek.Message{
					Role:    deepseek.ChatMessageRoleUser,
					Content: allInfo,
				})
				baseReq.Messages = dialog
				// tool调用没法流式处理
				if toolFlag {
					resp, err := cli.CreateChatCompletion(ctx, &baseReq)
					if err != nil {
						log.Fatalln(err)
					}
					msg := resp.Choices[0].Message
					dialog = append(dialog, msg)

					if len(msg.ToolCalls) == 0 {
						p.Cyan(msg.Content)
					} else {
						for _, toolCall := range msg.ToolCalls {
							// 从字典获取工具函数
							fn := toolCall.Function.Name
							toolCallRes := true
							if fn == "addFile" {
								tool.CreateFile(toolCall.Function.Arguments)
							} else if fn == "addDict" {
								tool.CreateDict(toolCall.Function.Arguments)
							} else if fn == "editFile" {
								tool.EditFile(toolCall.Function.Arguments)
							} else if fn == "isExist" {
								toolCallRes = tool.IsExist(toolCall.Function.Arguments)
							} else {
								log.Fatalln("fn not exist", fn, err)
							}

							// 添加工具响应
							dialog = append(dialog, deepseek.Message{
								Role:       deepseek.ChatMessageRoleTool,
								Content:    fmt.Sprintf("%v", toolCallRes),
								ToolCallID: toolCall.ID,
								Name:       toolCall.Function.Name,
							})
						}
						your()
						toolFlag = false
					}
				} else {
					baseReq.Stream = true
					resp, err := cli.CreateChatCompletionStream(ctx, baseReq)
					if err != nil {
						log.Fatalln(err)
					}
					go func() {
						var streamContent strings.Builder
						for {
							response, err := resp.Recv()
							if errors.Is(err, io.EOF) {
								// 流式响应完成，将完整内容加入对话历史
								dialog = append(dialog, deepseek.Message{
									Role:    deepseek.ChatMessageRoleAssistant,
									Content: streamContent.String(),
								})
								fmt.Println()
								your()
								return
							} else if err != nil {
								fmt.Printf("\nStream error: %v\n", err)
								return
							}

							if response.Choices[0].Delta.Content != "" {
								content := response.Choices[0].Delta.Content
								fmt.Print(content)
								streamContent.WriteString(content)
							}
						}
					}()
				}

				input.Reset()         // 清空输入缓冲区
				lastLineEmpty = false // 重置标记
			} else {
				lastLineEmpty = true
				continue
			}
		} else {
			input.WriteString(line + "\n")
			lastLineEmpty = false
		}
	}
}

func addFileToCtx(urls []string) {
	for _, url := range urls {
		is, info, err := tool.IsExistAndGetContent(url)
		if is && err == nil {
			dialog = append(dialog, deepseek.Message{
				Role:    deepseek.ChatMessageRoleUser,
				Content: fmt.Sprintf("url: %s \n,content: %s \n", url, info),
			})
		} else {
			if !is {
				log.Println("file isn't exist")
			} else {
				log.Println(info, is, err)
			}
		}
	}
	p.Yellow("add file to context success")
}

func addDictToCtx(urls []string) {
	for _, url := range urls {
		// 检查路径是否存在，并且是否为目录
		if _, err := os.Stat(url); err == nil {
			// 如果是目录，递归处理目录中的文件
			err := processDirectory(url)
			if err != nil {
				log.Println("Error processing directory:", err)
				continue
			}
		} else {
			log.Println("Directory doesn't exist:", url)
		}
	}
	p.Yellow("add directory to context success")
}

func processDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// 标记目录是否为空
	isEmpty := true

	for _, file := range files {
		filePath := dirPath + "/" + file.Name()
		if file.IsDir() {
			// 如果是子目录，递归处理
			err := processDirectory(filePath)
			if err != nil {
				return err
			}
			isEmpty = false // 子目录存在，目录不为空
		} else {
			// 如果是文件，调用 addFileToCtx
			addFileToCtx([]string{filePath})
			isEmpty = false // 文件存在，目录不为空
		}
	}

	// 如果目录为空，添加目录信息到上下文
	if isEmpty {
		addEmptyDirToCtx(dirPath)
	}

	return nil
}

func addEmptyDirToCtx(dirPath string) {
	fmt.Printf("add {empty Dir} [%s] to context\n", dirPath)
	dialog = append(dialog, deepseek.Message{
		Role:    deepseek.ChatMessageRoleUser,
		Content: fmt.Sprintf("url: %s \n [this is Empty Directory]\n", dirPath),
	})
}

func your() {
	p.Blue("Your input:")
}
func ai() {
	p.Cyan("ai output:")
}
