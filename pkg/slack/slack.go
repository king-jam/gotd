package slack

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"net/http"

	"github.com/king-jam/gotd/pkg/api/models"
	"github.com/king-jam/gotd/pkg/gif"
	"github.com/nlopes/slack"
)

const successMsg = "GIF Successfully posted to GOTD"

// New returns a handler for incoming Slack supported commands
func New(service gif.Service, verificationToken string) http.Handler {
	return slashCommandHandler{service: service, verificationToken: verificationToken}
}

type slashCommandHandler struct {
	service           gif.Service
	verificationToken string
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

	userCmd := "Requested GIF\n" + s.Text

	w.WriteHeader(http.StatusOK)

	err = h.updateGIF(s)
	if err != nil {
		log.Debugf("request failed: %s", err)
		msg := []byte(userCmd + "\n" + err.Error())

		if _, err := w.Write(msg); err != nil {
			log.Debugf("write failed: %s", err)
		}

		return
	}

	if _, err := w.Write([]byte(userCmd + "\n" + successMsg)); err != nil {
		log.Debugf("write failed with err: %s", err)
	}
}

const gotdCommand = "/gotd"

var errTokenInvalid = errors.New("unable to validate slack token")
var errInvalidCommand = errors.New("invalid slash command sent")
var errUnauthorizedUser = errors.New("user not authorized")
var errUpdating = errors.New("error setting")

func (h slashCommandHandler) updateGIF(s slack.SlashCommand) error {
	// Validate if command is from slack
	if !s.ValidateToken(h.verificationToken) {
		return errTokenInvalid
	}

	switch s.Command {
	case gotdCommand:
		// Validate user against the user pool
		if !h.validateUser(s.UserID) {
			return errUnauthorizedUser
		}

		newGif := &models.GIF{
			URL:         s.Text,
			RequestSrc:  "slack",
			RequesterID: s.UserID,
		}

		err := h.service.Set(context.Background(), newGif)
		if err != nil {
			return errUpdating
		}
	default:
		return errInvalidCommand
	}

	return nil
}
