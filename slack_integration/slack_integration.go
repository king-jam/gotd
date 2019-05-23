package slack_integration

import (
	"log"
	"net/http"
	"os"

	"github.com/king-jam/gotd/gif"
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
		// u, err := url.Parse(s.Text)
		// if err != nil {
		// 	w.WriteHeader(http.StatusOK)
		// 	w.Write([]byte(userCmd + "\n" + err.Error()))
		// 	return
		// }

		// err = normalizeGiphyURL(u)
		// if err != nil {
		// 	w.WriteHeader(http.StatusOK)
		// 	w.Write([]byte(userCmd + "\n" + err.Error()))
		// 	return
		// }
		// tags, err := giphy.GetGIFTags(u.String())
		// if err != nil {
		// 	w.WriteHeader(http.StatusOK)
		// 	w.Write([]byte(userCmd + "\n" + err.Error()))
		// 	return
		// }
		newGif := &gif.GIF{
			GIF:         s.Text,
			RequestSrc:  "slack",
			RequesterID: s.UserID,
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
