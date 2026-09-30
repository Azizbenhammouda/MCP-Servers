package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var scope string = "https://www.googleapis.com/auth/gmail.send"

const credentialsFile = "credentials.json"
const tokenFile = "token.json"

func main() {
	data, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Fatal(err)
	}

	cnf, err := google.ConfigFromJSON(data, scope)
	if err != nil {
		log.Fatal(err)
	}
	cnf.RedirectURL = "http://localhost:8080"

	codeCh := make(chan string)

	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprintln(w, "Success! You can close this tab and return to the terminal.")
		codeCh <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	authURL := cnf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Visit this URL, log in, and approve access:")
	fmt.Println(authURL)

	code := <-codeCh
	server.Shutdown(context.Background())

	tok, err := cnf.Exchange(context.Background(), code)
	if err != nil {
		log.Fatal(err)
	}

	tokJSON, err := json.MarshalIndent(tok, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(tokenFile, tokJSON, 0600); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nSaved token.json. Refresh token:")
	fmt.Println(tok.RefreshToken)
}
