package dashboard

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/king-jam/gotd/postgres"
)

func New(db *postgres.DBClient) http.Handler {
	return dashboardHandler{db: db}
}

type dashboardHandler struct {
	db *postgres.DBClient
}

func (d dashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gif, err := d.db.LatestGIF()
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
