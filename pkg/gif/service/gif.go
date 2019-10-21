// Package service provides the implementation of a
package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/king-jam/gotd/pkg/api/models"
	"github.com/king-jam/gotd/pkg/gif"
)

// Service provides
type Service struct {
	repo gif.Repository
}

// New returns an initialized service struct
func New(repo gif.Repository) *Service {
	return &Service{repo: repo}
}

// Set is the primary actor in the GIF service.
func (s *Service) Set(ctx context.Context, g *models.GIF) error {
	// and reformat the URL
	err := s.buildGIFFromURL(g)
	if err != nil {
		return errors.New("invalid URL")
	}
	//Update deactivate time for previous gif before storing new gif
	lastGif, err := s.repo.Last()
	if err != nil {
		// If there is no previous gif, then store new gif
		if err == gif.ErrRecordNotFound {
			err = s.repo.Insert(g)
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	// If user post the same URL twice, do nothing and return
	if g.URL == lastGif.URL {
		return nil
	}

	//Else, update the deactivate time for previous gif
	now := strfmt.DateTime(time.Now())
	lastGif.DeactivatedAt = &now

	err = s.repo.Update(lastGif)
	if err != nil {
		return err
	}
	// Insert gif into db
	err = s.repo.Insert(g)
	if err != nil {
		return err
	}

	return nil
}

// UpdateGif will update the current gif will necessary data
func (s *Service) UpdateGif(g *models.GIF) error {
	err := s.repo.Update(g)
	if err != nil {
		return err
	}

	return nil
}

// GetMostRecent returns the latest GIF that is active
func (s *Service) Latest(ctx context.Context) (*models.GIF, error) {
	g, err := s.repo.Last()
	if err != nil {
		return nil, err
	}

	return g, nil
}

// buildGIFFromURL parses and normalizes a request
func (s *Service) buildGIFFromURL(g *models.GIF) error {
	// Reformat the URL
	url, err := url.ParseRequestURI(g.URL)
	if err != nil {
		return err
	}

	// Normalize the URL
	err = normalizeGiphyURL(url)
	if err != nil {
		return err
	}

	g.URL = url.String()

	return nil
}

// validateURL will validate if URL is from giphy.com
func validateURL(url *url.URL) bool {
	// Validate if string is from giphy
	return url.Hostname() == "giphy.com"
}

// normalizeGiphyURL will add /fullscreen to URL
func normalizeGiphyURL(url *url.URL) error {
	if !validateURL(url) {
		return fmt.Errorf("invalid URL - Use Giphy.com")
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
