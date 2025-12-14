// agent/model.go
package agnet

import (
	"agentx/llm"
	"agentx/prompt"
	"agentx/tool"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Option func
type AgentOption func(AgentInterface)

// 设置子智能体
func WithSubAgents(AgentType AgentTypeOption, subAgents []AgentInterface) AgentOption {
	return func(a AgentInterface) {
		if agent, ok := a.(*Agent); ok {
			agent.AgentType = AgentType
			agent.SubAgents = subAgents
		}
	}
}

// 设置工具列表
func WithTools(tools ...tool.ToolInvoker) AgentOption {
	return func(a AgentInterface) {
		if agent, ok := a.(*Agent); ok {
			agent.Tools = tools
		}
	}
}

// 单个智能体结构
type Agent struct {
	Name      string             // 智能体名称
	Ctx       *context.Context   // 上下文
	Llm       llm.LlmInterface   // 使用的LLM模型
	Input     string             // 用户输入
	Prompt    *prompt.Prompt     // 提示词
	Output    string             // 智能体输出
	AgentType AgentTypeOption    // 智能体类型，单智能体，多智能体
	LoopMax   int                // 循环智能体最大循环次数
	Memory    bool               // 是否启用记忆功能
	SubAgents []AgentInterface   // 子智能体列表
	Tools     []tool.ToolInvoker // 智能体可用工具列表
}

// 初始化智能体
func NewAgent(ctx context.Context, name string, llmModel llm.LlmInterface, sysPrompt *prompt.Prompt, opts ...AgentOption) AgentInterface {
	model := &Agent{
		Name:      name,
		Llm:       llmModel,
		Prompt:    sysPrompt,
		AgentType: AgentTypeSingle,
		LoopMax:   5,
		Memory:    false,
	}
	for _, opt := range opts {
		opt(model)
	}
	return model
}

// 智能体名称
func (a *Agent) AgentName() string {
	return a.Name
}

func (a *Agent) buildToolPrompt() string {
	if len(a.Tools) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n你可以使用以下工具：\n")

	for _, t := range a.Tools {
		sb.WriteString(t.Name() + ":" + t.ToolPrompt())
		sb.WriteString("\n")
	}

	sb.WriteString(`
当你需要使用工具时，请严格输出 JSON：
{
  "tool": "工具名",
  "arguments": { ... }
}
如果不需要工具，请直接输出文本。
`)
	return sb.String()
}

// 执行智能体任务
// todo 每一种类型都单出一个函数拆分开来
func (a *Agent) Execute(input string) (string, error) {
	// 需要系统提示词的智能体类型（单智能体、并发父智能体）必须先初始化 Prompt
	if (a.AgentType == AgentTypeSingle || a.AgentType == AgentTypeParallel) && a.Prompt == nil {
		return "", errors.New("请先初始化系统提示语")
	}

	// todo 独立函数实现
	// 工作流执行（顺序执行）
	// 上一个智能体输出结果为当前智能体的输入
	if a.AgentType == AgentTypeSequential {
		var currentInput string = input
		var currentOutput string
		var err error
		for _, agent := range a.SubAgents {
			currentOutput, err = agent.Execute(currentInput)
			if err != nil {
				return "", err
			}

			// trace todo
			fmt.Println("========================")
			fmt.Printf("当前智能体[%v],当前输入[%v],当前输出[%v]\n", agent.AgentName(), currentInput, currentOutput)
			fmt.Println("========================")

			// 获取当前智能体的输出作为下一次输入
			currentInput = currentOutput
		}
		return currentOutput, nil
	}

	// todo 独立函数实现
	// 并发执行
	// 所有子智能体并发执行，最后父智能体汇总执行
	if a.AgentType == AgentTypeParallel {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var currentInput string = input
		var outputs []string

		// 并发执行子智能体
		for _, agent := range a.SubAgents {
			wg.Add(1)
			go func(agent AgentInterface) {
				defer wg.Done()
				output, err := agent.Execute(currentInput)
				// trace todo
				fmt.Println("========================")
				fmt.Printf("当前智能体[%v],当前输入[%v],当前输出[%v]\n", agent.AgentName(), currentInput, output)
				fmt.Println("========================")

				if err != nil {
					return
				}
				mu.Lock()
				outputs = append(outputs, output)
				mu.Unlock()
			}(agent)
		}
		wg.Wait()
		currentInput = strings.Join(outputs, "")

		// 父智能体执行
		res, err := a.Llm.Execute(a.Prompt.Prompt, currentInput)
		if err != nil {
			return "", err
		}
		a.Output = res

		return a.Output, nil
	}

	// 单智能体执行
	if a.Prompt == nil {
		return "", errors.New("请先初始化系统提示语")
	}

	// 1️⃣ 拼 Prompt（系统提示 + 工具说明）
	sysPrompt := a.Prompt.Prompt + a.buildToolPrompt()

	// 2️⃣ 第一次 LLM 推理
	resp, err := a.Llm.Execute(sysPrompt, input)
	if err != nil {
		return "", err
	}

	// 3️⃣ 判断是否 Tool Call
	// todo 改成 多工具调用格式 MCP
	tc, ok := parseToolCall(resp)
	if !ok {
		a.Output = resp
		return resp, nil
	}

	// 4️⃣ 执行 Tool
	// todo 改成多工具调用
	var result *tool.ToolResult
	for _, t := range a.Tools {
		if tc.Name == t.Name() {
			result, err = t.Call(tc)
			break
		}
	}
	if err != nil {
		return "", err
	}

	// 5️⃣ Tool 结果 → 再次喂给 LLM
	// todo 改成多工具调用
	finalResp, err := a.Llm.Execute(
		sysPrompt,
		fmt.Sprintf("工具调用结果如下：\n%s\n请基于结果回答用户。", result.Output),
	)
	if err != nil {
		return "", err
	}

	a.Output = finalResp
	return finalResp, nil
}

func parseToolCall(output string) (*tool.ToolCall, bool) {
	output = strings.TrimSpace(output)

	// todo 这种判断不精确，待优化
	if !strings.HasPrefix(output, "{") {
		return nil, false
	}

	var tc tool.ToolCall
	if err := json.Unmarshal([]byte(output), &tc); err != nil {
		return nil, false
	}
	if tc.Name == "" {
		return nil, false
	}
	return &tc, true
}
