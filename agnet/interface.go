package agnet

// 智能体基本能力
type AgentInterface interface {
	// 1. 基础能力
	// 执行智能体任务
	Execute(input string) (string, error)

	// 2. 多智能体能力
}
