package url

import (
	"errors"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

var (
	twitterRE    = regexp.MustCompile(`^https?://twitter\.com/.*?/status/\d+(/photo/\d+)?(#|$)`)
	tweetRE      = regexp.MustCompile(`(?s)<meta property="og:description" content="([^"]*)"`)
	tweetImageRE = regexp.MustCompile(`<meta property="og:image" content="([^"]*)"`)
)

func handleTwitter(target string) (string, error) {
	if !twitterRE.MatchString(target) {
		return "", errSkip
	}

	u, err := url.Parse(target)
	if err != nil {
		return "", err
	}
	// Twitter frontend is not scrapable easily and they require using their API.
	// Instead, use the alternative Twitter front-end https://github.com/zedeus/nitter
	u.Host = "nitter.net"

	body, err := get(u.String())
	if err != nil {
		return "", err
	}

	text := tweetRE.FindStringSubmatch(body)
	if text == nil {
		return "", errors.New("url: twitter: cannot parse tweet")
	}
	s := text[1]

	image := tweetImageRE.FindStringSubmatch(body)
	if image != nil {
		uri, err := url.QueryUnescape(image[1])
		if err != nil {
			return "", err
		}
		s = fmt.Sprintf("%v %v", s, uri)
	}

	s = stripTags(s)
	// unescape would replace &nbsp; by \u00a0 but we prefer normal space \u0020
	s = strings.Replace(s, "&nbsp;", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)
	s = html.UnescapeString(s)
	s = trim(s)
	return s, nil
}
