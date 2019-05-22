package dashboard

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/king-jam/gotd/gif"
)

func New(service *gif.GifService) http.Handler {
	return dashboardHandler{service: service}
}

type dashboardHandler struct {
	service *gif.GifService
}

func (d dashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gif, err := d.service.GetMostRecent()
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
