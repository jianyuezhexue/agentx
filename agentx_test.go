package agentx

import (
	"agentx/agnet"
	"agentx/llm"
	"agentx/prompt"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgent(t *testing.T) {
	ctx := context.Background()

	// 实例化提示词
	prompt := prompt.NewPrompt(nil, "你是一个智能助理，能够帮助用户完成各种任务。请根据用户的输入提供有用的信息和建议。")

	// 实例话千问大模型
	qwenModel := &llm.QwenWenModel{
		Token:        "sk-e692504205e74522b45710e1c25065ad",
		BaseUrl:      "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
		Model:        "qwen-plus",
		OutputFormat: llm.OutputFormatText,
		Temperature:  0.5,
		TopP:         0.5,
	}
	qwen := llm.NewQWenModel(qwenModel)

	// 初始化智能体
	demoAgent := agnet.NewAgent(ctx, qwen, prompt)
	res, err := demoAgent.Execute("帮我写一首关于春天的诗歌。")

	// 断言输出结果
	assert.Equal(t, nil, err)
	fmt.Println(res)
}
