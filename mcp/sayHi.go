// mcp/server/sayHi.go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	Name   string `json:"name" jsonschema:"the name of the person to greet"`
	Source string `json:"source,omitempty" jsonschema:"where this call comes from, e.g. 'direct' or 'agent'"`
}

type Output struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, input Input) (*mcp.CallToolResult, Output, error) {
	out := Output{Greeting: "Hi " + input.Name}
	src := input.Source
	if src == "" {
		// 默认认为是由 Agent 发起的工具调用
		src = "agent"
	}
	log.Printf("[MCP] SayHi called, source=%s, name=%s, greeting=%s", src, input.Name, out.Greeting)

	// 将输出内容追加写入本地 log.txt 文件
	// 文件格式：时间戳 + 空格 + [source] + 空格 + greeting 文本
	f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Printf("[MCP] open log.txt failed: %v", err)
		return nil, out, nil
	}
	defer f.Close()

	line := time.Now().Format(time.RFC3339) + " [" + src + "] " + out.Greeting + "\n"
	if _, err := f.WriteString(line); err != nil {
		log.Printf("[MCP] write log.txt failed: %v", err)
	}

	return nil, out, nil
}

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
