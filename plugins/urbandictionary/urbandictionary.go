// Package urbandictionary implements a plugin to get urban definition of words.
package urbandictionary

import (
	"encoding/json"
	"fmt"
	bot "github.com/StalkR/goircbot"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Result struct {
	Has_related_words bool
	Pages             int
	Result_type       string
	Sounds            []string
	Total             int
	List              []Definition
}

type Definition struct {
	Word                   string
	Definition             string
	Example                string
	Author                 string
	Defid                  int
	Permalink              string
	Current_vote           string
	Thumbs_up, Thumbs_down int
	Term                   string
	Type                   string
}

func (r *Result) String() string {
	if len(r.List) == 0 {
		return "no result"
	}
	if r.Result_type == "exact" || r.Result_type == "fulltext" {
		return r.List[0].String()
	}
	if r.Result_type == "no_results" {
		terms := make([]string, 0, len(r.List))
		for _, d := range r.List {
			terms = append(terms, d.Term)
		}
		return fmt.Sprintf("no result, did you mean: %s?", strings.Join(terms[:5], "? "))
	}
	return fmt.Sprintf("%s: %v", r.Result_type, r.List)
}

func (d *Definition) String() string {
	def := d.Definition
	def = strings.Replace(def, "\r", "", -1)
	def = strings.Replace(def, "\n", " ", -1)
	if len(def) > 200 {
		def = def[:200] + "..."
	}
	return fmt.Sprintf("%s: %s", d.Word, def)
}

// Define gets definition of term on urbandictionary and populates a Result.
func Define(term string) (r Result, e error) {
	base := "http://api.urbandictionary.com/v0/define"
	params := url.Values{}
	params.Set("term", term)
	resp, err := http.Get(fmt.Sprintf("%s?%s", base, params.Encode()))
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(contents, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func urban(b *bot.Bot, e *bot.Event) {
	term := strings.TrimSpace(e.Args)
	if len(term) == 0 {
		return
	}
	r, err := Define(term)
	if err != nil {
		b.Conn.Privmsg(e.Target, fmt.Sprintf("error: %s", err))
		return
	}
	b.Conn.Privmsg(e.Target, r.String())
}

// Register registers the plugin with a bot.
func Register(b *bot.Bot) {
	b.AddCommand("urban", bot.Command{
		Help:    "get definition of word from urbandictionary",
		Handler: urban,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
