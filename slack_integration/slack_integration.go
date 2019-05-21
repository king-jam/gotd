package slack_integration

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/king-jam/gotd/giphy"
	"github.com/king-jam/gotd/postgres"
	"github.com/nlopes/slack"
)

const SuccessMsg = "GIF Successfully posted to GOTD"

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
	userCmd := s.Command + " " + s.Text
	switch s.Command {
	case "/gotd":
		userId := s.UserID
		if !validateUser(userId) {
			response := userCmd + "\n" + "You don't have permission to change GOTD"
			w.Write([]byte(response))
			return
		}

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
		newGif := &postgres.GifHistory{
			GIF:         u.String(),
			RequestSrc:  "slack",
			RequesterID: s.UserID,
			Tags:        tags,
		}
		// Update deactivate time for previous gif
		lastGif, err := h.db.LatestGIF()
		if err != nil {
			if err == postgres.ErrRecordNotFound {
				err = h.db.Insert(newGif)
				if err != nil {
					log.Print(err)
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(userCmd + "\n" + err.Error()))
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(userCmd + "\n" + SuccessMsg))
				return
			}
			log.Print(err)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userCmd + "\n" + err.Error()))
			return
		}
		// if we are using the same exact GIF URL
		if lastGif.GIF == newGif.GIF {
			// just return success and change nothing
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userCmd + "\n" + SuccessMsg))
			return
		}
		lastGif.DeactivatedAt = time.Now()
		fmt.Printf("\n\n%+v\n\n", lastGif)
		err = h.db.Update(lastGif)
		if err != nil {
			log.Print("failed to update the last gif")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userCmd + "\n" + err.Error()))
			return
		}

		// Insert new gif into db
		err = h.db.Insert(newGif)
		if err != nil {
			log.Print("failed to insert into db")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userCmd + "\n" + err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userCmd + "\n" + SuccessMsg))
	default:
		msg := fmt.Errorf("invalid command")
		log.Print(msg)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userCmd + "\n" + msg.Error()))
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

func normalizeGiphyURL(url *url.URL) error {
	if !validateURL(url) {
		return fmt.Errorf("Invalid URL - Use Giphy.com")
	}
	var fullPath string
	// Check if URL has "/fullscreen"
	//basePath := path.Dir(url.Path)
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
