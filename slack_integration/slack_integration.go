package slack_integration

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/king-jam/gotd/gif"
	"github.com/king-jam/gotd/giphy"
	"github.com/nlopes/slack"
)

const SuccessMsg = "GIF Successfully posted to GOTD"

func New(service *gif.GifService) http.Handler {
	return slashCommandHandler{service: service}
}

type slashCommandHandler struct {
	service *gif.GifService
}

// Slack command handler
func (h slashCommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse the command
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Print("failed to parse command")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Validate if command is from slack
	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		log.Print("unable to validate Slack Token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userCmd := "Requested GIF\n" + s.Text
	switch s.Command {
	case "/gotd":
		// Validate user against the user pool
		userId := s.UserID
		if !validateUser(userId) {
			response := userCmd + "\n" + "You don't have permission to change GOTD"
			w.Write([]byte(response))
			return
		}

		// Grab the Gif URL from command
		u, err := url.Parse(s.Text)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userCmd + "\n" + err.Error()))
			return
		}

		err = normalizeGiphyURL(u)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userCmd + "\n" + err.Error()))
			return
		}
		tags, err := giphy.GetGIFTags(u.String())
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userCmd + "\n" + err.Error()))
			return
		}
		newGif := &gif.GIF{
			GIF:         u.String(),
			RequestSrc:  "slack",
			RequesterID: s.UserID,
			Tags:        tags,
		}

		err = h.service.StoreGif(newGif)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userCmd + "\n" + err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userCmd + "\n" + SuccessMsg))
	default:
		return
	}
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
