package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func generateNonce() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func percentEncode(str string) string {
	return url.QueryEscape(str)

}

func generateSignature(baseStr, signkey string) string {
	h := hmac.New(sha1.New, []byte(signkey))
	h.Write([]byte(baseStr))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))

}

func postTweet(tweet string) {
	consumerKey := os.Getenv("TWITTER_API_KEY")
	consumerSecret := os.Getenv("TWITTER_API_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	oauthNonce := generateNonce()
	oauthTimeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	requestMethod := "POST"
	requestURL := "https://api.x.com/1.1/statuses/update.json"
	params := url.Values{}
	params.Add("status", tweet)
	params.Add("include_entities", "true")
	params.Add("oauth_consumer_key", consumerKey)
	params.Add("oauth_nonce", oauthNonce)
	params.Add("oauth_signature_method", "HMAC-SHA1")
	params.Add("oauth_timestamp", oauthTimeStamp)
	params.Add("oauth_token", accessToken)
	params.Add("oauth_version", "1.0")

	baseStr := requestMethod + "&" + percentEncode(requestURL) + "&" + percentEncode(params.Encode())

	signingKey := percentEncode(consumerSecret) + "&" + percentEncode(accessSecret)

	oauthSignature := generateSignature(baseStr, signingKey)
	authHeader := fmt.Sprintf(`OAuth oauth_consumer_key="%s", oauth_nonce="%s", oauth_signature="%s", oauth_signature_method="HMAC-SHA1",oauth_timestamp="%s", oauth_token="%s", oauth_version="1.0"`,
		percentEncode(consumerKey), percentEncode(oauthNonce), percentEncode(oauthSignature), percentEncode(oauthTimeStamp), percentEncode(accessToken))

	println(authHeader)
	client := &http.Client{}
	data := url.Values{}
	data.Add("status", tweet)

	req, err := http.NewRequest("POST", requestURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatal("Error creating request: ", err)
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making request: ", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading response: ", err)
	}

	fmt.Println("Response Status:", res.Status)
	fmt.Println("Response Body:", string(body))
}

func main() {
	godotenv.Load()
	postTweet("Avneet is Here!!")
}
