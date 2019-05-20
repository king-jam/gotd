package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/king-jam/gotd/postgres"
	"github.com/nlopes/slack"
)

var UserList = []string{"val", "king-jam", "minh.nguyen", "Aman"}
var DB *postgres.DBClient

func getenv(name string) string {
	env := os.Getenv(name)
	if env == "" {
		return ""
	}
	return env
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Print("failed to parse command")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	members, err := getUserList()
	log.Print("List of members \n")
	log.Print(members)
	switch s.Command {
	case "/gotd":
		userName := s.UserName
		if !validUser(userName) {
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
		log.Printf("%s", s.ChannelName)
		if validateURL(u) {
			newGif := &postgres.GOTD{
				GIF: u.String(),
			}
			err := DB.UpdateGIF(newGif)
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

func validUser(userId string) bool {
	for _, user := range UserList {
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

func getUserList() ([]string, error) {
	channelId := os.Getenv("CHANNEL_ID")
	api := slack.New(os.Getenv("CLIENT_TOKEN"))
	channel, err := api.GetChannelInfo(channelId)
	if err != nil {
		log.Print("Error getting channel info")
		return nil, err
	}
	members := channel.Members
	log.Printf("Number of members in this channel: %d", channel.NumMembers)
	return members, nil
}

func gifHandler(w http.ResponseWriter, r *http.Request) {
	gif, err := DB.LatestGIF()
	if err != nil {
		log.Print("failed to get latest GIF")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(gif)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dbString := os.Getenv("DATABASE_URL")
	if dbString == "" {
		log.Fatal("$DATABASE_URL must be set")
	}

	dbURL, err := url.Parse(dbString)
	if err != nil {
		log.Fatal("Invalid Database URL format")
	}

	DB, err = postgres.InitDatabase(dbURL)
	if err != nil {
		log.Fatal("Unable to initialize the Database")
	}
	defer DB.Close()

	http.HandleFunc("/receive", slashCommandHandler)
	http.HandleFunc("/gif", gifHandler)
	http.Handle("/", http.FileServer(http.Dir("./static/dashboard")))

	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
