// tool/mcp_tool.go
package tool

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPTool struct {
	ToolName string
	Desc     string
	Session  *mcp.ClientSession
}

// MCPCalled 用于测试场景下确认是否真的走过 MCPTool.Call。
// 生产逻辑不会依赖这个标记。
var MCPCalled atomic.Bool

func (m *MCPTool) Name() string {
	return m.ToolName
}

func (m *MCPTool) ToolPrompt() string {
	return m.Desc
}

func (m *MCPTool) Call(tc *ToolCall) (*ToolResult, error) {
	if m.Session == nil {
		return nil, errors.New("mcp session is nil")
	}

	// 标记：已经通过 MCP 客户端发起工具调用（用于测试验证）。
	MCPCalled.Store(true)

	res, err := m.Session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      tc.Name,
		Arguments: tc.Arguments,
	})
	if err != nil {
		return nil, err
	}
	if res.IsError {
		return nil, fmt.Errorf("tool call error")
	}

	for _, c := range res.Content {
		if t, ok := c.(*mcp.TextContent); ok {
			return &ToolResult{Output: t.Text}, nil
		}
	}
	return &ToolResult{Output: ""}, nil
}
