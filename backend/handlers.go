package backend

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"encoding/json"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

// Method that process the request of a google notification to our server
func (s *Server) GoogleWatcherHandler(w http.ResponseWriter, r *http.Request) {
	if s.worker.IsClosed() {
		serverError(w)
		return
	}
	switch r.Method {
	case http.MethodPost:
		header := r.Header
		channelID := header.Get("X-Goog-Channel-ID")
		resourceID := header.Get("X-Goog-Resource-ID")
		resourceState := header.Get("X-Goog-Resource-State")
		log.Debugf("resourceID: %s from channelID: %s", resourceID, channelID)
		if resourceState == "sync" {
			w.WriteHeader(http.StatusOK)
			return
		}

		err := s.manageSynchronizationGoogle(channelID)
		var status int
		if err != nil {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusOK
		}
		w.WriteHeader(status)

	default:
		notFound(w)
		return
	}
}

// Method that process the request of a outlook notification to our server
func (s *Server) OutlookWatcherHandler(w http.ResponseWriter, r *http.Request) {
	if s.worker.IsClosed() {
		serverError(w)
		return
	}
	switch r.Method {
	case http.MethodPost:
		validationToken := r.FormValue("validationtoken")
		if len(validationToken) > 0 {
			log.Debugf("ValidationToken: %s", validationToken)
			w.Header().Set("Content-Type", "plain/text")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s", validationToken)
		} else {
			contents, err := ioutil.ReadAll(r.Body)
			if err != nil {
				serverError(w)
				return
			}
			notification := new(api.OutlookNotification)
			err = json.Unmarshal(contents, &notification)
			if err != nil {
				serverError(w)
				return
			}
			log.Warningf("OUTLOOK SUB: %s", contents)
			err = s.manageSynchronizationOutlook(notification.Subscriptions)
			var status int
			if err != nil {
				status = http.StatusInternalServerError
			} else {
				status = http.StatusOK
			}
			w.WriteHeader(status)
		}
		return
	default:
		notFound(w)
		return
	}

}

func notFound(w http.ResponseWriter) {
	contents, err := ioutil.ReadFile("./frontend/resources/html/404.html")
	if err != nil {
		serverError(w)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(contents)
}

func serverError(w http.ResponseWriter) {
	contents, err := ioutil.ReadFile("./frontend/resources/html/500.html")
	if err != nil {
		// TODO: something more elegant
		panic(err)
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(contents)

}
