package dashboard

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/king-jam/gotd/gif"
)

// New returns an initialized handler for the dashboard
func New(service *gif.Service) http.Handler {
	return dashboardHandler{service: service}
}

type dashboardHandler struct {
	service *gif.Service
}

func (d dashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gif, err := d.service.GetMostRecent()
	if err != nil {
		log.Debugf("failed to get latest GIF: %s", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(gif)
	if err != nil {
		log.Debugf("failed to marshal response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(response); err != nil {
		log.Debugf("response failed with err: %s", err)
	}
}
