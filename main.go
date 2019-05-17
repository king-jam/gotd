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

	switch s.Command {
	case "/gotd":
		userName := s.UserName
		if !validUser(userName) {
			response := "You are not a member of this channel"
			w.Write([]byte(response))
			return
		}
		channelID := s.ChannelID
		members, _ := getUserList(channelID)

		log.Printf("List of users %s", members)
		url := fmt.Sprintf("%v", slack.Msg{Text: s.Text}.Text)
		log.Print(url)
		newGif := &postgres.GOTD{
			GIF: url,
		}
		err := DB.UpdateGIF(newGif)
		if err != nil {
			log.Print("failed to insert into db")
			w.WriteHeader(http.StatusInternalServerError)
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

// func getUserList() ([]string, error) {
// 	token := os.Getenv("SLACK_VERIFICATION_TOKEN")
// 	api := slack.New(token)
// 	users, err := api.GetUsers()
// 	var usernames []string
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i := range users {
// 		usernames = append(usernames, users[i].Name)
// 	}
// 	return usernames, nil
// }

func getUserList(channelID string) ([]string, error) {
	token := os.Getenv("SLACK_VERIFICATION_TOKEN")
	api := slack.New(token)
	channel, err := api.GetChannelInfo(channelID)
	if err != nil {
		return nil, err
	}
	members := channel.Members
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
