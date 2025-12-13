package llm

import (
	"fmt"

	"resty.dev/v3"
)

// 千问大模型
type QwenWenModel struct {
	Token        string             // 千问大模型授权码
	BaseUrl      string             // 千问大模型服务地址
	Model        string             // 使用的模型名称
	Temperature  float64            // 采样温度
	TopP         float64            // 采样TopP
	OutputFormat OutputFormatOption // text, json
}

// 实例话千问大模型
func NewQWenModel(model *QwenWenModel) LlmInterface {
	return model
}

// 执行任务
func (qwm *QwenWenModel) Execute(sysPrompt string, input string) (string, error) {

	// 空值处理
	if input == "" {
		return "", nil
	}

	// 实例化HTTP客户端
	http := resty.New()
	defer http.Close()

	// 构建请求体
	requestBody := QwenReq{
		Model: qwm.Model,
		Messages: []*QwenMessage{
			{
				Role:    "system",
				Content: sysPrompt,
			},
			{
				Role:    "user",
				Content: input,
			},
		},
		ResponseFormat: &QwenResponseFormat{
			Type: "text",
		},
	}

	// 指定json返回类型
	if qwm.OutputFormat == OutputFormatJSON {
		requestBody.ResponseFormat.Type = "json_object"
	}

	// 设置请求头
	http.Header().Set("Authorization", fmt.Sprintf("Bearer %s", qwm.Token))
	http.Header().Set("Content-Type", "application/json")

	// 发送请求
	resp := &QwenResp{}
	_, err := http.R().SetBody(requestBody).SetResult(resp).Post(qwm.BaseUrl)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("error: %s", resp.Error.Message)
	}

	// 返回内容
	return resp.Choices[0].Message.Content, nil
}
