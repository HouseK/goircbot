// Package search implements Google Custom Search.
package search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

type Result struct {
	Kind              string
	Url               Url
	Queries           Queries
	Context           Context
	SearchInformation SearchInformation
	Items             []Item
}

type Url struct {
	Type, Template string
}

type Queries struct {
	NextPage, Request []Page
}

type Page struct {
	Title, TotalResults, SearchTerms string
	Count, StartIndex                int
	InputEncoding, OutputEncoding    string
	Safe, Cx                         string
}

type Context struct {
	Title string
}

type SearchInformation struct {
	SearchTime            float64
	FormattedSearchTime   string
	TotalResults          string
	FormattedTotalResults string
}

type Item struct {
	Kind, Title, HtmlTitle         string
	Link, DisplayLink              string
	Snippet, HtmlSnippet, CacheId  string
	FormattedUrl, HtmlFormattedUrl string
}

func compactSpaces(s string) string {
	r, err := regexp.Compile("\\s\\s+")
	if err != nil {
		return s
	}
	return string(r.ReplaceAll([]byte(s), []byte(" ")))
}

func (i *Item) String() string {
	return fmt.Sprintf("%s - %s", i.Link, compactSpaces(i.Snippet))
}

// Search searches a term on Google Custom Search and returns a Result.
// It requires a Google API Key (key) and a Google Custom Search ID (cx).
func Search(term, key, cx string) (*Result, error) {
	base := "https://www.googleapis.com/customsearch/v1"
	params := url.Values{}
	params.Set("key", key)
	params.Set("cx", cx)
	params.Set("alt", "json")
	params.Set("q", term)
	resp, err := http.Get(fmt.Sprintf("%s?%s", base, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := Result{}
	err = json.Unmarshal(contents, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
