package giphy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/html"
)

// Scrape for tags from given URL string
func GetGIFTags(gif string) ([]string, error) {
	tags := make([]string, 0, 0)
	var ok bool
	url, err := url.Parse(gif)
	if err != nil {
		fmt.Printf("ERROR: URL parsing failed: %s\n", err)
		return tags, err
	}
	trimFullScreenFromURL(url)
	resp, err := http.Get(url.String())
	if err != nil {
		fmt.Println("ERROR: Failed to get page")
		return tags, err
	}
	b := resp.Body
	defer b.Close() // close Body when the function returns
	z := html.NewTokenizer(b)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			cleanedTags := removeUnusedTags(tags)
			return cleanedTags, nil
		case html.StartTagToken:
			t := z.Token()
			// Check if the token is an <meta> tag
			isAnchor := t.Data == "meta"
			if !isAnchor {
				continue
			}

			// Determine if this is the keywords meta data
			if !isKeywordsToken(t) {
				continue
			}
			// Since we know this is the SEO keywords, extract the keywords aka tags
			ok, tags = getTags(t)
			if !ok {
				continue
			}
		}
	}
}

func isKeywordsToken(t html.Token) (ok bool) {
	// Iterate over all of the Token's attributes and check for "name" = "keywords"
	for _, a := range t.Attr {
		isTagAnchor := a.Key == "name" && a.Val == "keywords"
		if isTagAnchor {
			ok = true
		}
	}
	return
}

// Remove "GIF" and "Animated GIF" tags that come with when scraping
func removeUnusedTags(tags []string) []string {
	var reducedTagsLen int
	log.Printf("len(tags): %d", len(tags))
	if len(tags) < 3 {
		reducedTagsLen = len(tags)
	} else {
		reducedTagsLen = len(tags) - 2
	}
	reducedTags := make([]string, 0, reducedTagsLen)
	for _, t := range tags {
		if t != "GIF" && t != "Animated GIF" {
			reducedTags = append(reducedTags, t)
		}
	}
	return reducedTags
}

// getTags will get the gif's tags from token attribute
func getTags(t html.Token) (ok bool, tags []string) {
	// Iterate over all of the Token's attributes until we find an "content"
	for _, a := range t.Attr {
		if a.Key == "content" {
			tag := a.Val
			tags = strings.Split(tag, ", ")
			ok = true
		}
	}
	return
}

// trimFullScreenFromURL will remove "/fullscreen" from URL
func trimFullScreenFromURL(url *url.URL) {
	basePath := path.Dir(url.Path)
	url.Path = basePath
	return
}
