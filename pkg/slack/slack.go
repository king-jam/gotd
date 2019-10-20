package slack

import (
	log "github.com/sirupsen/logrus"

	"net/http"
	"os"

	"github.com/king-jam/gotd/pkg/gif"
	"github.com/nlopes/slack"
)

const successMsg = "GIF Successfully posted to GOTD"

// New returns a handler for incoming Slack supported commands
func New(service *gif.Service) http.Handler {
	return slashCommandHandler{service: service}
}

type slashCommandHandler struct {
	service *gif.Service
}

// Slack command handler
func (h slashCommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse the command
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Debugf("failed to parse command: %s", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	// Validate if command is from slack
	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		log.Debugf("unable to validate slack token: %s", err)
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	userCmd := "Requested GIF\n" + s.Text

	switch s.Command {
	case "/gotd":
		// Validate user against the user pool
		userID := s.UserID
		if !h.validateUser(userID) {
			response := userCmd + "\n" + "You don't have permission to change GOTD"
			if _, err := w.Write([]byte(response)); err != nil {
				log.Debugf("write failed with err: %s", err)
			}

			return
		}

		newGif := &gif.GIF{
			URL:         s.Text,
			RequestSrc:  "slack",
			RequesterID: s.UserID,
		}

		err = h.service.StoreGif(newGif)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusOK)

			if _, err := w.Write([]byte(userCmd + "\n" + err.Error())); err != nil {
				log.Debugf("write failed with err: %s", err)
			}

			return
		}

		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte(userCmd + "\n" + successMsg)); err != nil {
			log.Debugf("write failed with err: %s", err)
		}
	default:
		return
	}
}
