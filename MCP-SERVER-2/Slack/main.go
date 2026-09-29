package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// requires manual invitation of the bots into the private channel
type SlackInput struct {
	Channel string `json:"channel" jsonschema:"Slack channel where the message will be diffused"`
	Message string `json:"message" jsonschema:"the full notification text to send"`
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}

var token = requireEnv("SLACK_BOT_TOKEN")

func send_slack_notification(ctx context.Context, req *mcp.CallToolRequest, args SlackInput) (*mcp.CallToolResult, any, error) {
	err := SendSlackMessage(token, args.Channel, args.Message)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
		}, nil, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Slack message sent !"},
		},
	}, nil, nil
}
func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "SlackServer",
		Version: "1.0.0",
	}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_slack_notification",
		Description: "Send a Slack notification message to a given channel",
	}, send_slack_notification)
	err := server.Run(context.Background(), &mcp.StdioTransport{})
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
