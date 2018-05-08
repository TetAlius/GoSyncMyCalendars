package backend

import (
	"fmt"
	"net/http"
	"strings"

	"io/ioutil"

	"encoding/json"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

func (s *Server) GoogleTokenHandler(w http.ResponseWriter, r *http.Request) {
	route, err := util.CallAPIRoot("google/token/uri")
	log.Debugln(route)
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		serverError(w)
		return
	}
	params, err := util.CallAPIRoot("google/token/request-params")
	log.Debugln(params)
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		serverError(w)
		return
	}
	query := r.URL.Query()
	// TODO: Know how to send state
	//state := query.Get("state")

	code := query.Get("code")

	client := http.Client{}
	req, err := http.NewRequest("POST",
		route,
		strings.NewReader(
			fmt.Sprintf(params, code)))

	if err != nil {
		log.Errorf("error creating new google request: %s", err.Error())
		serverError(w)
		return
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error doing google request: %s", err.Error())
		serverError(w)
		return
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading response body from google request: %s", err.Error())
		serverError(w)
		return
	}

	//TODO: DB to implement
	account, err := api.NewGoogleAccount(contents)

	go func(account *api.GoogleAccount) {
		log.Debugln(account)
		account.GetAllCalendars()
		account.Refresh()
	}(account)

	//This is so that users cannot read the response
	http.Redirect(w, r, "http://localhost:8080", http.StatusPermanentRedirect)
}

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
		}
		//TODO: look here what was the change of the resource
		//Google does not give the change of resource
		//Possible changes include the creation of a new resource, or the modification or deletion of an existing resource.
		//channelID := header.Get("X-Goog-Channel-ID")
		//token := header.Get("X-Goog-Channel-Token")
		//expiration := header.Get("X-Goog-Channel-Expiration")
		resourceID := header.Get("X-Goog-Resource-ID")
		//resourceURI := header.Get("X-Goog-Resource-URI")
		//messageNumber := header.Get("X-Goog-Message-Number")
		//TODO: manage to get event with this id
		event := &api.GoogleEvent{ID: resourceID}
		log.Debugf("id of event to synchronize: %s", resourceID)
		err := event.SetState(api.UpdatedText)
		if err != nil {
			s.worker.Events <- event
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

	default:
		notFound(w)
		return
	}
}

func (s *Server) OutlookTokenHandler(w http.ResponseWriter, r *http.Request) {
	route, err := util.CallAPIRoot("outlook/token/uri")
	log.Debugln(route)
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		serverError(w)
		return
	}
	params, err := util.CallAPIRoot("outlook/token/request-params")
	log.Debugln(params)
	if err != nil {
		log.Errorf("error generating URL: %s", err.Error())
		serverError(w)
		return
	}

	client := http.Client{}
	code := r.URL.Query().Get("code")

	req, err := http.NewRequest("POST",
		route,
		strings.NewReader(
			fmt.Sprintf(params, code)))
	if err != nil {
		log.Errorf("error creating new outlook request: %s", err.Error())
		serverError(w)
		return
	}

	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error doing outlook request: %s", err.Error())
		serverError(w)
		return
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading response body from outlook request: %s", err.Error())
		serverError(w)
		return
	}
	log.Debugln(contents)
	//TODO: DB to implement
	account, err := api.NewOutlookAccount(contents)
	if err != nil {
		log.Errorf("error creating new account request: %s", err.Error())
		serverError(w)
		return
	}
	go func(account *api.OutlookAccount) {
		log.Debugln(account)
		account.GetAllCalendars()
		account.Refresh()
	}(account)

	http.Redirect(w, r, "http://localhost:8080", http.StatusPermanentRedirect)
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
			done := make(chan bool)
			go func() {
				for _, subs := range notification.Subscriptions {
					eventID := subs.Data.ID
					//TODO: manage to get event with this id
					event := &api.OutlookEvent{ID: eventID}
					log.Debugf("id of event to synchronize: %s", eventID)
					err := event.SetState(subs.ChangeType)
					if err != nil {
						s.worker.Events <- event
					} else {
						done <- false
					}
				}
				done <- true
			}()
			processed := <-done
			if processed {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
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
