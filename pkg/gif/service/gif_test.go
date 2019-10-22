// Package service provides the implementation of a
package service

import (
	"net/url"
	"strings"
	"testing"

	"github.com/king-jam/gotd/pkg/api/models"
)

func TestService_buildGIFFromURL(t *testing.T) {
	tests := map[string]struct {
		url     string
		wantErr bool
	}{
		"good no fullscreen": {"https://giphy.com/gifs/patriots-gopats-ZCe9VWewCrNBEpTqFs", false},
		"good fullscreen":    {"https://giphy.com/gifs/exit-gate-employee-2IIT1RVK0UCpq/fullscreen", false},
		"bad not giphy":      {"https://mysupersite.com", true},
		"not a URL":          {"BADURL", true},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			s := &Service{}
			g := &models.GIF{URL: tt.url}
			if err := s.buildGIFFromURL(g); (err != nil) != tt.wantErr {
				t.Errorf("Service.buildGIFFromURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_normalizeGiphyURL(t *testing.T) {
	tests := map[string]struct {
		url               string
		expectTransformed bool
		wantErr           bool
	}{
		"good no fullscreen": {"https://giphy.com/gifs/patriots-gopats-ZCe9VWewCrNBEpTqFs", true, false},
		"good fullscreen":    {"https://giphy.com/gifs/exit-gate-employee-2IIT1RVK0UCpq/fullscreen", false, false},
		"bad not giphy":      {"https://mysupersite.com", false, true},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			url, err := url.Parse(tt.url)
			if err != nil {
				t.Errorf("URL Fails to parse")
			}
			if err := normalizeGiphyURL(url); (err != nil) != tt.wantErr {
				t.Errorf("normalizeGiphyURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.expectTransformed {
				fsAdded := strings.Contains(url.Path, "fullscreen")
				if !fsAdded {
					t.Errorf("normalizeGiphyURL() transformed = %v, expectTransformed %v", fsAdded, tt.expectTransformed)
				}
			}
		})
	}
}
