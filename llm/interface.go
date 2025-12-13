package llm

// 智能体基本能力
type LlmInterface interface {
	// 执行LLM模型任务
	Execute(sysPrompt string, input string) (string, error)
}
