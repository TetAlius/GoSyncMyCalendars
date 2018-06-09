package backend

import (
	"net"
	"net/http"

	"fmt"

	"context"
	"time"

	"strings"

	"encoding/json"

	"database/sql"

	"os"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/backend/db"
	"github.com/TetAlius/GoSyncMyCalendars/convert"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/worker"
	"github.com/getsentry/raven-go"
)

//Backend object
type Server struct {
	IP       net.IP
	Port     int
	server   *http.Server
	mux      *http.ServeMux
	worker   *worker.Worker
	database db.Database
	ticker   *time.Ticker
	sentry   *raven.Client
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//NewBackend creates a backend
func NewServer(ip string, port int, maxWorker int, database *sql.DB, sentry *raven.Client) *Server {
	data := db.New(database, sentry)
	server := Server{IP: net.ParseIP(ip), Port: port, mux: http.NewServeMux(), worker: worker.New(maxWorker, data), database: data, sentry: sentry}
	server.mux.HandleFunc("/google/watcher", server.GoogleWatcherHandler)
	server.mux.HandleFunc("/outlook/watcher", server.OutlookWatcherHandler)
	server.mux.HandleFunc("/accounts/", server.retrieveInfoHandler)
	server.mux.HandleFunc("/subscribe/", server.subscribeCalendarHandler)
	server.mux.HandleFunc("/refresh/", server.refreshHandler)
	return &server
}

//Start the backend
func (s *Server) Start() (err error) {
	go s.worker.Start()
	log.Debugln("Start backend")

	laddr := fmt.Sprintf("%s:%d", s.IP.String(), s.Port)
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: s,
	}
	log.Infof("Backend server listening at %s", laddr)
	go s.manageSubscriptions()

	err = s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Errorf("ListenAndServe: " + err.Error())
	}
	return
}

//Stop the backend
func (s *Server) Stop() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	log.Debugf("Stopping backend with ctx: %s", ctx)
	var returnErr error
	err = s.server.Shutdown(ctx)
	if err != nil {
		s.sentry.CaptureErrorAndWait(err, map[string]string{"stopping": "backend server"})
		returnErr = err
	}
	err = s.worker.Stop()
	if err != nil {
		s.sentry.CaptureErrorAndWait(err, map[string]string{"stopping": "backend worker"})
		returnErr = err
	}
	s.ticker.Stop()
	err = s.database.Close()
	if err != nil {
		s.sentry.CaptureErrorAndWait(err, map[string]string{"stopping": "backend database"})
		returnErr = err
	}
	return returnErr
}

func (s *Server) subscribeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	ok := manageCORS(w, *r, map[string]bool{"POST": true, "DELETE": true})
	if !ok {
		return
	}
	email, userUUID, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//decoded := strings.Split(string(auth), ":")
	log.Debugf("email: %s", email)
	log.Debugf("user: %s", userUUID)
	if len(email) == 0 || len(userUUID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	param := r.URL.Path[len("/subscribe/"):]
	switch r.Method {
	case http.MethodPost:
		calendar, err := s.database.RetrieveCalendars(email, userUUID, param)
		if err != nil {
			log.Errorf("error getting calendar")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Debugf("%s", calendar)
		err = prepareSync(calendar)
		if err != nil {
			log.Errorf("error starting sync")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = s.database.StartSync(calendar, userUUID)
		if err != nil {
			log.Errorf("error trying to start sync: %s", calendar.GetUUID())
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodDelete:
		log.Debugf("Getting method delete")
		err := s.database.StopSync(param, email, userUUID)
		if err != nil {
			log.Errorf("error stopping sync: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

	}
}

func (s *Server) retrieveInfoHandler(w http.ResponseWriter, r *http.Request) {
	ok := manageCORS(w, *r, map[string]bool{"GET": true})
	if !ok {
		return
	}
	email, userUUID, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Debugf("email: %s", email)
	log.Debugf("user: %s", userUUID)
	if len(email) == 0 || len(userUUID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	param := r.URL.Path[len("/accounts/"):]
	values := strings.Split(param, "/")
	account, err := s.database.RetrieveAccount(userUUID, values[0])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}
	err = account.Refresh()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = s.database.UpdateAccountFromUser(account, userUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	calendars, err := account.GetAllCalendars()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	contents, err := json.Marshal(calendars)
	w.Write(contents)
}

func (s *Server) refreshHandler(w http.ResponseWriter, r *http.Request) {
	ok := manageCORS(w, *r, map[string]bool{"POST": true})
	if !ok {
		return
	}
	email, userUUID, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Debugf("email: %s", email)
	log.Debugf("user: %s", userUUID)
	if len(email) == 0 || len(userUUID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := s.database.UpdateAllCalendarsFromUser(userUUID, email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func manageCORS(w http.ResponseWriter, r http.Request, methods map[string]bool) (ok bool) {
	ok = true
	keys := make([]string, len(methods))
	for key := range methods {
		keys = append(keys, key)
	}

	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Access-Control-Allow-Methods", fmt.Sprintf("OPTIONS,%s", strings.Join(keys, ",")))
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return false
	}
	if !methods[r.Method] {
		log.Errorf("method not supported: %s", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return

}

func (s *Server) manageSubscriptions() {
	s.ticker = updateTicker()
	for {
		<-s.ticker.C
		log.Debugf("next ticking: %s")
		subscriptions, err := s.database.GetExpiredSubscriptions()
		if err != nil {
			log.Errorf("error: %s", err.Error())
			s.ticker = updateTicker()
			continue
		}
		for _, subscription := range subscriptions {
			acc := subscription.GetAccount()
			if err := acc.Refresh(); err != nil {
				continue
			}
			if err = s.database.UpdateAccountFromSubscription(acc, subscription); err != nil {
				log.Errorf("error updating account: %s", err.Error())
			}
			err = subscription.Renew()
			if err != nil {
				continue
			}
			err := s.database.UpdateSubscription(subscription)
			if err != nil {
				log.Errorf("error updating subscription: %s", err.Error())
			}
		}
		s.ticker = updateTicker()
	}
}

func updateTicker() *time.Ticker {
	tim := time.Now()
	nextTick := time.Date(tim.Year(), tim.Month(), tim.Day(), 0, 5, 0, 0, time.Local)
	if !nextTick.After(time.Now()) {
		nextTick = nextTick.Add(time.Hour * 24)
	}
	diff := nextTick.Sub(time.Now())
	log.Debugf("next tick: %s", nextTick)
	return time.NewTicker(diff)
}

func prepareSync(calendar api.CalendarManager) (err error) {
	err = calendar.GetAccount().Refresh()
	if err != nil {
		log.Errorf("error refreshing account: %s", err.Error())
		return
	}

	cal, err := calendar.GetAccount().GetCalendar(calendar.GetID())
	convert.Convert(cal, calendar)
	for _, calen := range calendar.GetCalendars() {
		err := convert.Convert(calendar, calen)
		if err != nil {
			log.Errorf("error converting info: %s", err.Error())
			return err
		}
		log.Debugf("Name1: %s Name2: %s", calendar.GetName(), calen.GetName())
		err = calen.GetAccount().Refresh()
		if err != nil {
			log.Errorf("error refreshing account calendar: %s error: %s", calen.GetID(), err.Error())
			return err
		}
		err = calen.Update()

		if err != nil {
			log.Errorf("error updating calendar: %s error: %s", calen.GetID(), err.Error())
			return err
		}
	}
	return
}
