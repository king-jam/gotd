package giphy

import (
	"net/url"

	libgiphy "github.com/sanzaru/go-giphy"
)

type GiphySource interface {
	SearchGifByTags([]string) (url.URL, error)
}

type Giphy struct {
	api *libgiphy.Giphy
}

func NewGiphy(api_key string) (*Giphy, error) {
	api := libgiphy.NewGiphy(api_key)
	return &Giphy{api: api}, nil
}

func (g *Giphy) SearchGifByTags(tags string) (string, error) {
	res, err := g.api.GetSearch(tags, 1, -1, "pg", "", false)
	if err != nil {
		return "", err
	}
	return res.Data[0].Url, nil
}
