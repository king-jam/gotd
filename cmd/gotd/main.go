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
	case "/gotd":
		gifURL := &slack.Msg{Text: s.Text}.Text
		userName := &slack.Msg{Text: s.Text}.Username
		isMember := validateUser(userName)

		if !isMember {
			return
		}
		err := db.Insert(db, GOTD{gifURL: gifURL})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func validateUser(userId string) bool {
	for i := range UserList {
		if userId == UserList[i] {
			return true
		}
	}
	return false
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

	db, err := postgres.InitDatabase(dbURL)
	if err != nil {
		log.Fatal("Unable to initialize the Database")
	}
	defer db.Close()

	http.HandleFunc("/receive", slashCommandHandler)

	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":8080", nil)
}
