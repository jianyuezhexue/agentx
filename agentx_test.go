package agentx

import (
	"agentx/agnet"
	"agentx/llm"
	"agentx/prompt"
	"agentx/tool"
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

// 测试单个智能体
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
	demoAgent := agnet.NewAgent(ctx, "诗歌智能体", qwen, prompt)
	res, err := demoAgent.Execute("帮我写一首关于春天的诗歌。")

	// 断言输出结果
	assert.Equal(t, nil, err)
	fmt.Println(res)
}

// 测试工作流智能体
func TestSequentialAgent(t *testing.T) {
	ctx := context.Background()

	// 实例话千问大模型
	qwenModel := &llm.QwenWenModel{
		Token:        "sk-e692504205e74522b45710e1c25065ad",
		BaseUrl:      "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
		Model:        "qwen-plus",
		OutputFormat: llm.OutputFormatText,
		Temperature:  0.5,
	}
	qwen := llm.NewQWenModel(qwenModel)

	// 智能体2
	prompt2 := prompt.NewPrompt(nil, "请将以下内容翻译成英文")
	agent2 := agnet.NewAgent(ctx, "翻译智能体", qwen, prompt2)

	// 智能体3
	prompt3 := prompt.NewPrompt(nil, "用用户输入的英文,写两日常句子,并注明中文翻译")
	agent3 := agnet.NewAgent(ctx, "诗歌智能体", qwen, prompt3)

	// 初始化工作流智能体
	withSubAgents := agnet.WithSubAgents(agnet.AgentTypeSequential, []agnet.AgentInterface{agent2, agent3})
	workflowAgent := agnet.NewAgent(ctx, "工作流智能体", qwen, nil, withSubAgents)

	// 执行工作流智能体任务
	res, err := workflowAgent.Execute("马拉松")

	// 断言输出结果
	assert.Equal(t, nil, err)
	fmt.Println(res)
}

// 测试并发智能体
func TestParallelAgent(t *testing.T) {
	ctx := context.Background()

	// 实例话千问大模型
	qwenModel := &llm.QwenWenModel{
		Token:        "sk-e692504205e74522b45710e1c25065ad",
		BaseUrl:      "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
		Model:        "qwen-plus",
		OutputFormat: llm.OutputFormatText,
		Temperature:  0.5,
	}
	qwen := llm.NewQWenModel(qwenModel)

	// 智能体2
	prompt2 := prompt.NewPrompt(nil, "将用户输入的哲理,引用经典著作进行中文解释")
	agent2 := agnet.NewAgent(ctx, "哲理解释中文", qwen, prompt2)

	// 智能体3
	prompt3 := prompt.NewPrompt(nil, "将用户输入的哲理,引用经典著作进行英文解释,返回英文内容")
	agent3 := agnet.NewAgent(ctx, "哲理解释英文", qwen, prompt3)

	// 初始化并发智能体
	withSubAgents := agnet.WithSubAgents(agnet.AgentTypeParallel, []agnet.AgentInterface{agent2, agent3})
	prompt0 := prompt.NewPrompt(nil, "将两个智能体返回的结果原样展示，不做任何翻译,并在前面追加一份简单总结，不超过30字。")
	parallelAgent := agnet.NewAgent(ctx, "并发智能体", qwen, prompt0, withSubAgents)

	// 执行并发智能体任务
	res, err := parallelAgent.Execute("真正的成熟，不是看清世界有多复杂，而是明白复杂之后依然选择善良与清醒。")

	// 断言输出结果
	assert.Equal(t, nil, err)
	fmt.Println(res)
}

// 测试Mcp调用agent
func TestMcp(t *testing.T) {
	ctx := context.Background()

	// 1️⃣ 启动 MCP Server（子进程）
	// 使用 `go run` 启动，避免对预编译二进制和可执行权限的依赖
	cmd := exec.Command("go", "run", "./mcp/sayHi.go")
	transport := &mcp.CommandTransport{Command: cmd}

	client := mcp.NewClient(
		&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"},
		nil,
	)
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Skipf("skipping TestMcp: failed to connect MCP server: %v", err)
	}
	defer session.Close()

	// 2️⃣ 构造 MCPTool
	sayHiTool := &tool.MCPTool{
		ToolName: "greet",
		Desc: `
工具名：greet
功能：向某人打招呼
参数：
- name(string)：名字
`,
		Session: session,
	}

	// 3️⃣ 构造真实 LLM
	qwen := llm.NewQWenModel(&llm.QwenWenModel{
		Token:        "sk-e692504205e74522b45710e1c25065ad",
		BaseUrl:      "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
		Model:        "qwen-plus",
		OutputFormat: llm.OutputFormatText,
		Temperature:  0.5,
	})

	// 4️⃣ 构造 Agent（注入 Tool）
	p := prompt.NewPrompt(nil, "你是一个智能助理，可以在必要时调用工具。")
	agent := agnet.NewAgent(ctx, "测试mcp智能体", qwen, p, agnet.WithTools(sayHiTool))

	// 5️⃣ 用“必须用工具”的输入
	out, err := agent.Execute("帮我向 Tom 打个招呼")
	assert.NoError(t, err)
	fmt.Println("Agent Output:", out)
}
