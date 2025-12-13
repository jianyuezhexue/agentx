package agnet

import (
	"agentx/llm"
	"agentx/prompt"
	"context"
)

// 单个智能体结构
type Agent struct {
	Ctx       *context.Context // 上下文
	Llm       llm.LlmInterface // 使用的LLM模型
	Input     string           // 用户输入
	Prompt    *prompt.Prompt   // 提示词
	Output    string           // 智能体输出
	AgentType AgentTypeOption  // 智能体类型，单智能体，多智能体
	LoopMax   int              // 循环智能体最大循环次数
	Memory    bool             // 是否启用记忆功能
	SubAgents []*Agent         // 子智能体列表
}

// 初始化智能体
func NewAgent(ctx context.Context, llmModel llm.LlmInterface, sysPrompt *prompt.Prompt) AgentInterface {
	return &Agent{
		Llm:       llmModel,
		Prompt:    sysPrompt,
		AgentType: AgentTypeSingle,
		LoopMax:   5,
		Memory:    false,
	}
}

// 执行智能体任务
func (a *Agent) Execute(input string) (string, error) {
	res, err := a.Llm.Execute(a.Prompt.Prompt, input)
	return res, err
}
