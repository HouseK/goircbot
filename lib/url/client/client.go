// Binary client shows the title of an url.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/StalkR/goircbot/lib/url"
)

var (
	flagURL          = flag.String("url", "", "URL to show title of.")
	flagTwitterToken = flag.String("twitter_token", "", "Twitter API Token.")
)

func main() {
	flag.Parse()
	if *flagURL == "" {
		fmt.Printf("Usage: %v [-twitter_token <token>] -url <url>\n", os.Args[0])
		os.Exit(1)
	}
	url.TwitterAPIToken = *flagTwitterToken
	title, err := url.Title(*flagURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(title)
}
