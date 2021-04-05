package main

import (
	arg "github.com/alexflint/go-arg"
	"github.com/cbrgm/clickbaiter/clickbaiter"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
	"runtime"
	"time"
)

var (
	// Version of clickbaiter-bot.
	Version string
	// Revision or Commit this binary was built from.
	Revision string
	// BuildDate this binary was built.
	BuildDate string
	// GoVersion running this binary.
	GoVersion = runtime.Version()
	// StartTime has the time this was started.
	StartTime = time.Now()
)

// Config represents command line flags.
type Config struct {
	TweetInterval         int    `arg:"env:TWEET_INTERVAL"`
	TwitterAccessToken    string `arg:"env:TWITTER_ACCESS_TOKEN"`
	TwitterSecretToken    string `arg:"env:TWITTER_SECRET_TOKEN"`
	TwitterConsumerKey    string `arg:"env:TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret string `arg:"env:TWITTER_CONSUMER_SECRET"`
}

func main() {
	conf := Config{
		TweetInterval: 60,
	}

	arg.MustParse(&conf)

	if conf.TwitterConsumerKey == "" || conf.TwitterConsumerSecret == "" || conf.TwitterAccessToken == "" || conf.TwitterSecretToken == "" {
		panic("Consumer key/secret and Access token/secret required")
	}

	filterOption := level.AllowInfo()

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, filterOption)
	logger = log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	level.Info(logger).Log(
		"msg", "starting clickbaiter-bot",
		"version", Version,
		"revision", Revision,
		"buildDate", BuildDate,
		"goVersion", GoVersion,
	)

	ticker := time.NewTicker(time.Duration(conf.TweetInterval) * time.Minute)
	publisher := NewPublisher(conf.TwitterConsumerKey, conf.TwitterConsumerSecret, conf.TwitterAccessToken, conf.TwitterSecretToken)

	cbg := clickbaiter.NewGenerator()
	cbg.UseHashtags(true)

	for {
		select {
		case <-ticker.C:
			err := publisher.PublishTweet(cbg.RandomSentence())
			if err != nil {
				level.Error(logger).Log("err", err)
			}
		}
	}
}

// Client represents the twitter client
type TwitterPublisher struct {
	client *twitter.Client
}

// NewFromConfig returns a new twitter Client from accountConfig
func NewPublisher(ConsumerKey string, ConsumerSecret string, AccessToken string, AccessSecret string) *TwitterPublisher {
	config := oauth1.NewConfig(ConsumerKey, ConsumerSecret)
	token := oauth1.NewToken(AccessToken, AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	return &TwitterPublisher{
		client: client,
	}
}

// PublishTweet published a new tweet with a given title and url as content
func (publisher *TwitterPublisher) PublishTweet(msg string) error {
	_, _, err := publisher.client.Statuses.Update(msg, nil)
	if err != nil {
		return err
	}
	return nil
}
