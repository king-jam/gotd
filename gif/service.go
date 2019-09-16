package gif

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/king-jam/gotd/giphy"
)

type GifService struct {
	repo  *Repo
	giphy *giphy.Giphy
}

func NewGifService(repo *Repo, giphy *giphy.Giphy) *GifService {
	return &GifService{repo: repo, giphy: giphy}
}

func (g *GifService) BuildGifFromUrl(gif *GIF) error {
	// Reformat the URL
	url, err := url.Parse(gif.URL)
	if err != nil {
		return err
	}

	//Get Tags From Gif URL
	tags, err := giphy.GetGIFTags(gif.URL)
	if err != nil {
		return err
	}

	// Normalize the URL
	err = normalizeGiphyURL(url)
	if err != nil {
		return err
	}
	gif.URL = url.String()
	var tagStructs []Tag
	for _, tag := range tags {
		t := Tag{
			Value: tag,
		}
		tagStructs = append(tagStructs, t)
	}
	gif.Tags = tagStructs
	return nil
}

func (g *GifService) BuildGifFromTags(gif *GIF) error {
	var b strings.Builder
	for _, t := range gif.Tags {
		b.WriteString(t.Value)
		b.WriteString(" ")
	}
	query := strings.TrimSpace(b.String())
	res, err := g.giphy.SearchGifByTags(query)
	if err != nil {
		return err
	}
	url, err := url.Parse(res)
	if err != nil {
		return err
	}
	err = normalizeGiphyURL(url)
	if err != nil {
		return err
	}
	gif.URL = url.String()
	return nil
}

func (g *GifService) StoreGif(gif *GIF) error {
	// Add more details onto the gif, such as tags, and reformat the URL
	validURL := govalidator.IsURL(gif.URL)
	if validURL {
		err := g.BuildGifFromUrl(gif)
		if err != nil {
			return err
		}
	} else {
		err := g.BuildGifFromTags(gif)
		if err != nil {
			return err
		}
	}
	//Update deactive time for previous gif before storing new gif
	lastGif, err := g.GetMostRecent()
	if err != nil {
		// If there is no previous gif, then store new gif
		if err == ErrRecordNotFound {
			err = g.repo.Insert(gif)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// If user post the same URL twice, do nothing and return
	if gif.URL == lastGif.URL {
		return nil
	}

	//Else, update the deactivate time for previous gif
	lastGif.DeactivatedAt = time.Now()
	err = g.UpdateGif(lastGif)
	if err != nil {
		return err
	}
	// Insert gif into db
	err = g.repo.Insert(gif)
	if err != nil {
		return err
	}
	return nil
}

func (g *GifService) UpdateGif(gif *GIF) error {
	err := g.repo.Update(gif)
	if err != nil {
		return err
	}
	return nil
}

func (g *GifService) RemoveGifById(id int) error {
	err := g.repo.DeleteGIFByID(id)
	if err != nil {
		return err
	}
	return nil
}

func (g *GifService) GetGIFById(id uint) (*GIF, error) {
	gif, err := g.repo.FindGIFByID(id)
	if err != nil {
		return &GIF{}, err
	}
	//object := TransformDBGifToGif(gif)
	return gif, nil
}

func (g *GifService) GetAllGifs() ([]GIF, error) {
	var gifList []GIF
	gifs, err := g.repo.FindAllGifs()
	if err != nil {
		return gifList, err
	}
	for _, g := range gifs {
		//gif := TransformDBGifToGif(&g)
		gifList = append(gifList, g)
	}
	return gifList, nil
}

func (g *GifService) GetMostRecent() (*GIF, error) {
	gif, err := g.repo.LatestGIF()
	if err != nil {
		return &GIF{}, err
	}
	//gif := TransformDBGifToGif(dbGif)
	return gif, nil
}

// validateURL will validate if URL is from giphy.com
func validateURL(url *url.URL) bool {
	// Validate if string is from giphy
	return url.Hostname() == "giphy.com"
}

// normalizeGiphyURL will add /fullscreen to URL
func normalizeGiphyURL(url *url.URL) error {
	if !validateURL(url) {
		return fmt.Errorf("Invalid URL - Use Giphy.com")
	}
	var fullPath string
	// Check if URL has "/fullscreen"
	ok, err := path.Match("/gifs/*/fullscreen", url.Path)
	if err != nil {
		return err
	}
	if !ok {
		fullPath = path.Join(url.Path, "fullscreen")
		url.Path = fullPath
	}
	return nil
}
