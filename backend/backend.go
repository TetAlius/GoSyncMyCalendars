package backend

import (
	"net"
	"net/http"

	"fmt"

	"context"
	"time"

	"strings"

	"encoding/base64"
	"encoding/json"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/backend/db"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/worker"
)

//Backend object
type Server struct {
	IP     net.IP
	Port   int
	server *http.Server
	mux    *http.ServeMux
	worker *worker.Worker
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	if cookie != nil {
		log.Debugf("session", cookie.Value)
		ctx := context.WithValue(r.Context(), "Session", cookie.Value)
		s.mux.ServeHTTP(w, r.WithContext(ctx))
	} else {
		s.mux.ServeHTTP(w, r)
	}
}

//NewBackend creates a backend
func NewServer(ip string, port int, maxWorker int) *Server {
	server := Server{IP: net.ParseIP(ip), Port: port, mux: http.NewServeMux(), worker: worker.New(maxWorker)}
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
	h := &http.Server{Addr: fmt.Sprintf(":%d", s.Port), Handler: s}
	s.server = h
	log.Infof("Backend server listening at %s", laddr)

	err = s.server.ListenAndServe()
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}
	return
}

//Stop the backend
func (s *Server) Stop() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	log.Debugf("Stopping backend with ctx: %s", ctx)
	err = s.server.Shutdown(ctx)
	s.worker.Stop()
	return
}

func (s *Server) subscribeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	ok := manageCORS(w, *r, map[string]bool{"POST": true, "DELETE": true})
	if !ok {
		return
	}
	authorization := r.Header.Get("Authorization")
	auth, err := base64.StdEncoding.DecodeString(authorization)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	decoded := strings.Split(string(auth), ":")
	log.Debugf("email: %s", decoded[0])
	log.Debugf("user: %s", decoded[1])
	if len(decoded[0]) == 0 || len(decoded[1]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	param := r.URL.Path[len("/subscribe/"):]
	calendar, err := db.RetrieveCalendars(decoded[0], decoded[1], param)
	if err != nil {
		log.Errorf("error getting calendar")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodPost:
		log.Debugf("%s", calendar)
		err = api.StartSync(calendar)
		if err != nil {
			log.Errorf("error starting sync")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		db.UpdateAccountFromUser(calendar.GetAccount(), decoded[1])
		db.UpdateCalendarFromUser(calendar, decoded[1])
		var subs api.SubscriptionManager
		switch calendar.(type) {
		case *api.GoogleCalendar:
			//TODO: change this IDS
			subs = api.NewGoogleSubscription("123", "URL")
			subs.Subscribe(calendar)
		case *api.OutlookCalendar:
			subs = api.NewOutlookSubscription("234", "URL", "Created,Deleted,Updated")
			subs.Subscribe(calendar)
		}
		db.SaveSubscription(subs, calendar)
		for _, cal := range calendar.GetCalendars() {
			db.UpdateAccountFromUser(cal.GetAccount(), decoded[1])
			db.UpdateCalendarFromUser(cal, decoded[1])
			var subscript api.SubscriptionManager
			switch cal.(type) {
			case *api.GoogleCalendar:
				subscript = api.NewGoogleSubscription("345", "URL")
				subscript.Subscribe(cal)
			case *api.OutlookCalendar:
				subscript = api.NewOutlookSubscription("456", "URL", "Created,Deleted,Updated")
				subscript.Subscribe(cal)
			}
			db.SaveSubscription(subscript, cal)
		}
	case http.MethodDelete:
		//	TODO: Delete all subscriptions and stoping them

	}
}

func (s *Server) retrieveInfoHandler(w http.ResponseWriter, r *http.Request) {
	ok := manageCORS(w, *r, map[string]bool{"GET": true})
	if !ok {
		return
	}
	authorization := r.Header.Get("Authorization")
	auth, err := base64.StdEncoding.DecodeString(authorization)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	decoded := strings.Split(string(auth), ":")
	log.Debugf("email: %s", decoded[0])
	log.Debugf("user: %s", decoded[1])
	if len(decoded[0]) == 0 || len(decoded[1]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	param := r.URL.Path[len("/accounts/"):]
	values := strings.Split(param, "/")
	account, err := db.RetrieveAccount(decoded[1], values[0])
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
	err = db.UpdateAccountFromUser(account, decoded[1])
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
	authorization := r.Header.Get("Authorization")
	log.Debugf("AUTH %s", authorization)
	auth, err := base64.StdEncoding.DecodeString(authorization)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	decoded := strings.Split(string(auth), ":")
	log.Debugln(decoded)
	log.Debugf("email: %s", decoded[0])
	log.Debugf("user: %s", decoded[1])
	if len(decoded[0]) == 0 || len(decoded[1]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = db.UpdateAllCalendarsFromUser(decoded[1], decoded[0])
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

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
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
