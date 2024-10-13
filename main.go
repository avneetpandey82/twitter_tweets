package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dghubble/oauth1"
)

func postTweet(client *http.Client, tweet string) {
	urlStr := "https://api.twitter.com/2/tweets"
	payload := url.Values{}
	payload.Set("status", tweet)
	res, err := client.PostForm(urlStr, payload)
	if err != nil {
		log.Fatalf("Error posting tweet: %v", err)
	}
	defer res.Body.Close()
	fmt.Println("Response Status: ", res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	fmt.Println("Tweet posted successfully: ", string(body))
}

func main() {
	consumerKey := os.Getenv("TWITTER_API_KEY")
	consumerSecret := os.Getenv("TWITTER_API_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	postTweet(httpClient, "Avneet is Here!!")
}
