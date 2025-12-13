package llm

// ollama模型
type OllamaModel struct {
	Model   string // 使用的模型名称
	BaseUrl string // ollama服务地址
}

// 实例化ollama模型
func NewOllamaModel(baseUrl, model string) LlmInterface {
	return &OllamaModel{
		BaseUrl: baseUrl,
		Model:   model,
	}
}

// 执行任务
func (om *OllamaModel) Execute(prompt string, input string) (string, error) {
	return "", nil
}
