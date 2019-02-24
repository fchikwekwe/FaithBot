package main

import (
	"fmt"
	"log"
	"os"
	// Import go-twitter modules
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Credentials stores all of our access/consumer tokens and secret keys needed
// for authentication against the twitter REST API.
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// getClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
func getClient(creds *Credentials) (*twitter.Client, error) {
	// Pass in the consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in the Access Token and the Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("User's Account:\n%+v\n", user)
	return client, nil
}

func sendTweet(client *twitter.Client) *twitter.Tweet {
	tweet, resp, err := client.Statuses.Update("A test tweet from a new bot I'm building!", nil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", resp)
	log.Printf("%+v\n", tweet)
	return tweet
}

func searchTweets(client *twitter.Client) *twitter.Search {
	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: "Golang",
	})
	if err != nil {
		log.Print(err)
	}
	// log.Printf("%+v\n", resp)
	log.Println("\n\n", search.Statuses[0].Text, search.Statuses[0].ID)
	return search
}

func sendRetweet(client *twitter.Client) {
	search := searchTweets(client)
	retweet, _, err := client.Statuses.Retweet(search.Statuses[0].ID, &twitter.StatusRetweetParams{
		ID: search.Statuses[0].ID,
	})
	if err != nil {
		log.Print(err)
	}
	// log.Printf("%+v\n", resp)
	log.Printf("%+v\n", retweet)
}

func likeTweet(client *twitter.Client) {
	search := searchTweets(client)
	like, _, err := client.Favorites.Create(&twitter.FavoriteCreateParams{
		ID: search.Statuses[0].ID,
	})
	if err != nil {
		log.Print(err)
	}
	log.Printf("%+v\n", like)
}

func main() {
	fmt.Println("Go-Twitter Bot v0.02")
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	client, err := getClient(&creds)
	// fmt.Println("client", client, creds)

	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
	}

	// fmt.Println("TYPE", reflect.TypeOf(client))
	// sendTweet(client)
	// searchTweets(client)
	// sendRetweet(client)
	likeTweet(client)

}
