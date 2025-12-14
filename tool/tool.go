package tool

type ToolCall struct {
	Name      string                 `json:"tool"`
	Arguments map[string]interface{} `json:"arguments"`
}

type ToolResult struct {
	Output string
}

type ToolInvoker interface {
	Name() string
	Call(tc *ToolCall) (*ToolResult, error)
	ToolPrompt() string // 返回给 LLM 用的工具说明
}
