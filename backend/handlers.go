package backend

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"encoding/json"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func (s *Server) GoogleWatcherHandler(w http.ResponseWriter, r *http.Request) {
	if s.worker.IsClosed() {
		serverError(w)
		return
	}
	switch r.Method {
	case http.MethodPost:
		header := r.Header
		resourceState := header.Get("X-Goog-Resource-State")
		if resourceState == "sync" {
			w.WriteHeader(http.StatusOK)
			return
		}
		//TODO: look here what was the change of the resource
		//Google does not give the change of resource
		//Possible changes include the creation of a new resource, or the modification or deletion of an existing resource.
		channelID := header.Get("X-Goog-Channel-ID")
		//token := header.Get("X-Goog-Channel-Token")
		//expiration := header.Get("X-Goog-Channel-Expiration")
		resourceID := header.Get("X-Goog-Resource-ID")
		//resourceURI := header.Get("X-Goog-Resource-URI")
		//messageNumber := header.Get("X-Goog-Message-Number")
		err := s.manageSynchronizationGoogle(channelID, resourceID)
		var status int
		if err != nil {
			status = http.StatusOK
		} else {
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)

	default:
		notFound(w)
		return
	}
}

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
			err = s.manageSynchronizationOutlook(notification.Subscriptions)
			var status int
			if err != nil {
				status = http.StatusOK
			} else {
				status = http.StatusInternalServerError
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
