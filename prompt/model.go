package prompt

import "context"

// options
type PromptOptions func(*Prompt)

// withPromptRule 设置任务规则
func WithPromptRule(rule string) PromptOptions {
	return func(p *Prompt) {
		p.PromptRule = rule
	}
}

// withPromptExample 设置任务示例
func WithPromptExample(example string) PromptOptions {
	return func(p *Prompt) {
		p.PromptExample = example
	}
}

// withBackgroundInfo 设置背景资料
func WithBackgroundInfo(info string) PromptOptions {
	return func(p *Prompt) {
		p.BackgroundInfo = info
	}
}

// withOutputConvention 设置输出约定
func WithOutputConvention(convention string) PromptOptions {
	return func(p *Prompt) {
		p.OutputConvention = convention
	}
}

// 结构化提示词
type Prompt struct {
	Ctx              *context.Context // 智能体上下文
	Prompt           string           // 系统提示词
	PromptRule       string           // 任务规则,我要什么,不要什么
	PromptExample    string           // 任务示例
	BackgroundInfo   string           // 背景资料
	OutputConvention string           // 输出约定
	MemorySchema     string           // 记忆schema
}

func NewPrompt(ctx *context.Context, prompt string, optipns ...PromptOptions) *Prompt {
	model := &Prompt{
		Ctx:    ctx,
		Prompt: prompt,
	}

	for _, item := range optipns {
		item(model)
	}

	return model
}

// 初始化提示词
func (p *Prompt) InitPrompt() string {
	if p.PromptRule != "" {
		p.Prompt += "\n任务规则:\n" + p.PromptRule
	}
	if p.PromptExample != "" {
		p.Prompt += "\n任务示例:\n" + p.PromptExample
	}
	if p.BackgroundInfo != "" {
		p.Prompt += "\n背景资料:\n" + p.BackgroundInfo
	}
	if p.OutputConvention != "" {
		p.Prompt += "\n输出约定:\n" + p.OutputConvention
	}
	return p.Prompt
}
