package agnet

// 智能体类型选项 ｜ 单智能体(chatModel)，工作流(sequential),并发(parallel),循环(loop)
type AgentTypeOption string

const (
	// 单智能体
	AgentTypeSingle AgentTypeOption = "single"
	// 工作流
	AgentTypeSequential AgentTypeOption = "sequential"
	// 并发
	AgentTypeParallel AgentTypeOption = "parallel"
	// 循环
	AgentTypeLoop AgentTypeOption = "loop"
)
