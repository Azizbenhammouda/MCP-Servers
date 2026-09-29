package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SlackMessagePayload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type SlackResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func SendSlackMessage(token, channel, message string) error {
	var response SlackResponse
	url := "https://slack.com/api/chat.postMessage"
	payload := SlackMessagePayload{
		Channel: channel,
		Text:    message,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	tkob := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", tkob)
	req.Header.Set("Content-Type", "application/json")
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
	err = json.Unmarshal(newData, &response)
	if err != nil {
		return err
	}
	if !response.Ok {
		return fmt.Errorf("slack api error: %s", response.Error)
	}
	return nil
}
