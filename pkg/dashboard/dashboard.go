package dashboard

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/king-jam/gotd/pkg/gif"
)

// New returns an initialized handler for the dashboard
func New(service gif.Service) http.Handler {
	cache := new(sync.Map)
	return dashboardHandler{service: service, cacheMap: cache}
}

type dashboardHandler struct {
	service  gif.Service
	cacheMap *sync.Map
}

func (d dashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gif, err := d.service.Latest(context.Background())
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
