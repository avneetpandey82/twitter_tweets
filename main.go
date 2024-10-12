package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/clientcredentials"
)

type TweetResponse struct {
	Data struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
}

func OAuth() *http.Client {
	config := clientcredentials.Config{
		ClientID:     os.Getenv("TWITTER_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	return config.Client(context.Background())
}

var baseUrl = `https://api.x.com/2`

func createPost(text string) {
	client := OAuth()
	tweetData := map[string]string{
		"text": text,
	}
	body, err := json.Marshal(tweetData)
	if err != nil {
		log.Fatal("Error marshaling tweet data: %v", err)
	}
	req, err := http.NewRequest("POST", baseUrl+`/tweets`, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TWITTER_BEARER_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making Post request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		log.Fatal("Error creating tweet: %v", res.StatusCode)
	}
	var tweetResponse TweetResponse
	if err := json.NewDecoder(res.Body).Decode(&tweetResponse); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Printf("Tweet created: %v (ID: %v)\n", tweetResponse.Data.Text, tweetResponse.Data.ID)

}

func main() {
	err := godotenv.Load()
	createPost("Hello! How are you?")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}
