package giphy

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/html"
)

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

func extractTag(url string) string {
	s := strings.Split(url, "/")
	return s[len(s)-1]
}

func isKeywordsToken(t html.Token) (ok bool) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		isTagAnchor := a.Key == "name" && a.Val == "keywords"
		if isTagAnchor {
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

func removeUnusedTags(tags []string) []string {
	reducedTags := make([]string, 0, len(tags)-2)
	for _, t := range tags {
		if t != "GIF" && t != "Animated GIF" {
			reducedTags = append(reducedTags, t)
		}
	}
	return reducedTags
}

func getTags(t html.Token) (ok bool, tags []string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "content" {
			tag := a.Val
			tags = strings.Split(tag, ", ")
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

func trimFullScreenFromURL(url *url.URL) {
	basePath := path.Dir(url.Path)
	url.Path = basePath
	return
}
