package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/king-jam/gotd/postgres"
	"github.com/nlopes/slack"
)

var UserList = []string{"val", "kingj2"}
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "gotd":
		userName := slack.Msg{Text: s.Text}.Username
		if !validUser(userName) {
			return
		}
		newGif := postgres.GOTD{
			GIF: slack.Msg{Text: s.Text}.Text,
		}
		err := DB.Insert(newGif)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func validUser(userId string) bool {
	return true
	// for i := range UserList {
	// 	if userId == UserList[i] {
	// 		return true
	// 	}
	// }
	// return false
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

	DB, err := postgres.InitDatabase(dbURL)
	if err != nil {
		log.Fatal("Unable to initialize the Database")
	}
	defer DB.Close()

	http.HandleFunc("/receive", slashCommandHandler)

	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
