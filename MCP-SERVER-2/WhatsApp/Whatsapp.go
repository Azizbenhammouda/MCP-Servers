package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ParameterDef struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ComponentDef struct {
	Type       string         `json:"type"`
	Parameters []ParameterDef `json:"parameters"`
}

type TemplateDef struct {
	Name       string            `json:"name"`
	Language   map[string]string `json:"language"`
	Components []ComponentDef    `json:"components"`
}

type MessagePayload struct {
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Template         TemplateDef `json:"template"`
}

func SendWhatsAppMessage(token, phoneNumberID, templateName, to, message string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v21.0/%s/messages", phoneNumberID)
	sent := MessagePayload{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: TemplateDef{
			Name: templateName,
			Language: map[string]string{
				"code": "en_US",
			},
			Components: []ComponentDef{
				{
					Type: "body",
					Parameters: []ParameterDef{
						{
							Type: "text",
							Text: message,
						},
					},
				},
			},
		},
	}
	data, err := json.Marshal(sent)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	tkob := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tkob)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	newData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return fmt.Errorf("whatsapp api error (status %d): %s", res.StatusCode, string(newData))
	}

	return nil
}
