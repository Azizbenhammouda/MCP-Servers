package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GmailInput struct {
	To      string `json:"to" jsonschema:"recipient email address"`
	Subject string `json:"subject" jsonschema:"email subject line"`
	Message string `json:"message" jsonschema:"the full plain-text email body"`
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}

var clientID = requireEnv("GMAIL_CLIENT_ID")
var clientSecret = requireEnv("GMAIL_CLIENT_SECRET")
var refreshToken = requireEnv("GMAIL_REFRESH_TOKEN")

func send_gmail_notification(ctx context.Context, req *mcp.CallToolRequest, args GmailInput) (*mcp.CallToolResult, any, error) {
	err := SendGmailMessage(clientID, clientSecret, refreshToken, args.To, args.Subject, args.Message)
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
			&mcp.TextContent{Text: "Email sent!"},
		},
	}, nil, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "GmailServer",
		Version: "1.0.0",
	}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_gmail_notification",
		Description: "Send a plain-text email notification to a given address",
	}, send_gmail_notification)
	err := server.Run(context.Background(), &mcp.StdioTransport{})
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
