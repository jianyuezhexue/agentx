package llm

// 千问大模型
type QwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QwenResponseFormat struct {
	Type string `json:"type"`
}

// 千问请求
type QwenReq struct {
	Model          string              `json:"model"`
	Messages       []*QwenMessage      `json:"messages"`
	ResponseFormat *QwenResponseFormat `json:"response_format"`
	// 采样温度，控制模型生成文本的多样性。temperature越高，生成的文本更多样，反之，生成的文本更确定。取值范围： [0, 2)
	Temperature float64 `json:"temperature"` // [0, 2)
	// 核采样的概率阈值，控制模型生成文本的多样性。top_p越高，生成的文本更多样。反之，生成的文本更确定。取值范围：（0,1.0]。
	TopP float64 `json:"top_p"` // （0,1.0]
	// 使用混合思考模型时，是否开启思考模式，适用于 Qwen3 、Qwen3-VL模型
	EnableThinking bool `json:"enable_thinking"`
}

// 千问返回

type Error struct {
	Code    string `json:"code"`
	Param   string `json:"param"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type Choice struct {
	Message      ResMessage   `json:"message"`
	FinishReason string       `json:"finish_reason"`
	Index        int          `json:"index"`
	Logprobs     *interface{} `json:"logprobs,omitempty"` // 使用 interface{} 因为 logprobs 的值是 null
}

type ResMessage struct {
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content"`
	Role             string `json:"role"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
type QwenResp struct {
	Choices           []*Choice `json:"choices"`
	Object            string    `json:"object"`
	Usage             Usage     `json:"usage"`
	Created           int64     `json:"created"`
	SystemFingerprint *string   `json:"system_fingerprint,omitempty"`
	Model             string    `json:"model"`
	ID                string    `json:"id"`
	Error             *Error    `json:"error,omitempty"`
}
