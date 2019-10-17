package gif

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/asaskevich/govalidator"
)

// Service provides 
type Service struct {
	repo  *Repo
}

// NewGifService returns an initialized service struct
func NewGifService(repo *Repo) *Service {
	return &Service{repo: repo}
}

// BuildGIFFromURL parses and normalizes a request
func (g *Service) BuildGIFFromURL(gif *GIF) error {
	// Reformat the URL
	url, err := url.Parse(gif.URL)
	if err != nil {
		return err
	}

	// Normalize the URL
	err = normalizeGiphyURL(url)
	if err != nil {
		return err
	}
	gif.URL = url.String()

	return nil
}

// StoreGif is the primary actor in the GIF service.
func (g *Service) StoreGif(gif *GIF) error {
	// Add more details onto the gif, such as tags, and reformat the URL
	validURL := govalidator.IsURL(gif.URL)
	if validURL {
		err := g.BuildGIFFromURL(gif)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Invalid URL")
		// err := g.BuildGifFromTags(gif)
		// if err != nil {
		// 	return err
		// }
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
	now := time.Now()
	lastGif.DeactivatedAt = &now
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

// UpdateGif will update the current gif will necessary data
func (g *Service) UpdateGif(gif *GIF) error {
	err := g.repo.Update(gif)
	if err != nil {
		return err
	}
	return nil
}

// GetMostRecent returns the latest GIF that is active
func (g *Service) GetMostRecent() (*GIF, error) {
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
