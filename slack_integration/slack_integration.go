package slack_integration

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/king-jam/gotd/postgres"
	"github.com/nlopes/slack"
)

var UserIdList = []string{
	"U5SFY08HW", // Ethan
	"U5SFZ590Q", // Val
	"UGG0Y2W82", //Aman
	"U5UAGKX4L", //Amy
	"U5U133V3Q", //Geoff
	"U5U0X61DM", // Joe
	"U5U1DSEQ7", // Justin
	"U61HFJ7V2", // Kranti
	"UFJRQ2S2F", // Minh
	"UFDAJLGJU", // Viet
	"U5V5T2DPZ", // Dale
	"UGYDW6UJK", // Edgardo
	"U5T9HLMAN", // James King
	"UEK11RZJP", // Sammie
	"UHH0LLBND", // Sandesh
	"U5SFZ590Q", // Val
}

func New(db *postgres.DBClient) http.Handler {
	return slashCommandHandler{db: db}
}

type slashCommandHandler struct {
	db *postgres.DBClient
}

func (h slashCommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Print("failed to parse command")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		log.Print("unable to validate Slack Token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/gotd":
		userId := s.UserID
		if !validateUser(userId) {
			response := "You don't have permission to change GOTD"
			w.Write([]byte(response))
			return
		}

		u, err := url.Parse(s.Text)
		if err != nil {
			response := "Invalid URL provided"
			w.WriteHeader(http.StatusPreconditionFailed)
			w.Write([]byte(response))
			return
		}
		if validateURL(u) {
			newGif := &postgres.CurrentGOTD{
				GIF: u.String(),
			}
			err := h.db.UpdateGIF(newGif)
			if err != nil {
				log.Print("failed to insert into db")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			response := "Please use giphy for your gif"
			w.Write([]byte(response))
			return
		}
	default:
		log.Print("invalid command")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func validateUser(userId string) bool {
	for _, user := range UserIdList {
		if userId == user {
			return true
		}
	}
	return false
}

func validateURL(url *url.URL) bool {
	// Validate if string is from giphy
	return url.Hostname() == "giphy.com"
}
