package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
}

func getAccessToken(clientID, clientSecret, refreshToken string) (string, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	res, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", err
	}
	if tr.AccessToken == "" {
		return "", fmt.Errorf("failed to get access token: %s", string(body))
	}
	return tr.AccessToken, nil
}

func SendGmailMessage(clientID, clientSecret, refreshToken, to, subject, message string) error {
	accessToken, err := getAccessToken(clientID, clientSecret, refreshToken)
	if err != nil {
		return err
	}

	rawEmail := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s",
		to, subject, message)

	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(rawEmail))

	payload := map[string]string{"raw": encoded}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://gmail.googleapis.com/gmail/v1/users/me/messages/send", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("gmail api error (status %d): %s", res.StatusCode, string(body))
	}
	return nil
}
