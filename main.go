package main

import (
	"log"
	"os"
	// Import gorm and sqlite
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

	saveTweet(tweet.ID, tweet.Text, "tweet")

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
	saveTweet(search.Statuses[0].ID, search.Statuses[0].Text, "retweet")

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
	saveTweet(search.Statuses[0].ID, search.Statuses[0].Text, "like")

	log.Printf("%+v\n", like)
}

// tweet is a struct that saves tweets for the user. It records the tweet ID,
// the tweet's text and the action that was taken by the user on that tweet
type tweet struct {
	gorm.Model
	tweetID int64
	Text    string
	Action  string // options: like, tweet, retweet
}

func saveTweet(tweetID int64, tweetText string, tweetAction string) {
	// Open DB
	db, err := gorm.Open("postgres", "host=http://faithbot.herokuapp.com/ port=7822")
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	// Create
	db.Create(&tweet{tweetID: tweetID, Text: tweetText, Action: tweetAction})
}

func main() {
	// Get auth credentials and the Twitter client
	creds := GetCreds()
	client, err := GetClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
	}
	// searchQuery := "Golang"
	testTweet := "*beep* Another test tweet from my bot. Checking data persistence. *beep*"
	// Examples of how to use the various functions this bot has
	SendTweet(client, testTweet)
	// SearchTweets(client, searchQuery)
	// SendRetweet(client, searchQuery)
	// LikeTweet(client, searchQuery)
}
