package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SendNotificationInput struct {
	To      string `json:"to" jsonschema:"recipient phone number in E.164 format, e.g. +911234567890"`
	Message string `json:"message" jsonschema:"the full notification text to send"`
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}

var token = requireEnv("WHATSAPP_TOKEN")
var phoneNumberID = requireEnv("WHATSAPP_PHONE_NUMBER_ID")
var templateName = requireEnv("WHATSAPP_TEMPLATE_NAME")

func send_whatsapp_notification(ctx context.Context, req *mcp.CallToolRequest, args SendNotificationInput) (*mcp.CallToolResult, any, error) {
	err := SendWhatsAppMessage(token, phoneNumberID, templateName, args.To, args.Message)
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
			&mcp.TextContent{Text: "Message sent !"},
		},
	}, nil, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "WhatsAppServer",
		Version: "1.0.0",
	}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_whatsapp_notification",
		Description: "Send a WhatsApp notification message to a given phone number",
	}, send_whatsapp_notification)
	err := server.Run(context.Background(), &mcp.StdioTransport{})
	if err != nil {
		log.Fatalf("Server failure %v", err)
	}
}
