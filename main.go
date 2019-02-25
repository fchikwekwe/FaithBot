package main

import (
	"database/sql"
	"log"
	"os"
	// Import mySQL

	_ "github.com/go-sql-driver/mysql"

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

// GetClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
func GetClient(creds *Credentials) (*twitter.Client, error) {
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

// GetCreds creates and returns a struct that gets the required environment
// variables for user authentication with the Twitter API
func GetCreds() Credentials {
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}
	return creds
}

// SendTweet sends a tweet with the specified text passed in as a string and
// returns a Tweet object.
func SendTweet(client *twitter.Client, tweetText string) *twitter.Tweet {
	tweet, _, err := client.Statuses.Update(tweetText, nil)
	if err != nil {
		log.Println(err)
	}
	// log.Printf("%+v\n", resp)
	log.Printf("%+v\n", tweet)
	return tweet
}

// SearchTweets searches for the given hashtag. It takes hashtag as the query
// argument (type string) as well as the twitter client. It returns a slice of
// Tweet objects.
func SearchTweets(client *twitter.Client, query string) *twitter.Search {
	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: query,
	})
	if err != nil {
		log.Print(err)
	}
	// log.Printf("%+v\n", resp)
	log.Println("\n\n", search.Statuses[0].Text, search.Statuses[0].ID)
	return search
}

// SendRetweet retweets the first retruned tweet after searching with the given
// hashtag. The hashtag must be passed in as a string along with the twitter
// client.
func SendRetweet(client *twitter.Client, searchQuery string) {
	search := SearchTweets(client, searchQuery)
	retweet, _, err := client.Statuses.Retweet(search.Statuses[0].ID, &twitter.StatusRetweetParams{
		ID: search.Statuses[0].ID,
	})
	if err != nil {
		log.Print(err)
	}
	// log.Printf("%+v\n", resp)
	log.Printf("%+v\n", retweet)
}

// LikeTweet sends a like to the first returned tweet after searching with the
// given hashtag. The hashtag is passed in as a string along with the twitter
// client.
func LikeTweet(client *twitter.Client, searchQuery string) {
	search := SearchTweets(client, searchQuery)
	like, _, err := client.Favorites.Create(&twitter.FavoriteCreateParams{
		ID: search.Statuses[0].ID,
	})
	if err != nil {
		log.Print(err)
	}
	log.Printf("%+v\n", like)
}

type tweet struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Action string `json:"action"`
}

func accessDB() {
	dbPass := os.Getenv("DB_PASS")
	db, err := sql.Open("mysql", "root:"+dbPass+"@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("CREATE DATABASE tweetRecall")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("USE tweetRecall")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	insert, err := db.Query("INSERT INTO test VALUES (2, 'TEST' )")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}
func main() {
	// Get auth credentials and the Twitter client
	// creds := GetCreds()
	// client, err := GetClient(&creds)
	// if err != nil {
	// 	log.Println("Error getting Twitter Client")
	// 	log.Println(err)
	// }
	// searchQuery := "Golang"
	// testTweet := "*beep* Test tweet from my bot. *beep*"
	// Examples of how to use the various functions this bot has
	// SendTweet(client, testTweet)
	// SearchTweets(client, searchQuery)
	// SendRetweet(client, searchQuery)
	// LikeTweet(client, searchQuery)
}
