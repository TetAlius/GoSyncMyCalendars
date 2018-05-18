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

func (s *Server) retrieveInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}
	authorization := r.Header.Get("Authorization")
	auth, err := base64.StdEncoding.DecodeString(authorization)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	decoded := strings.Split(string(auth), ":")
	log.Debugf("email: %s", decoded[0])
	log.Debugf("user: %s", decoded[1])
	if len(decoded[0]) == 0 || len(decoded[1]) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	param := r.URL.Path[len("/accounts/"):]
	values := strings.Split(param, "/")
	account, err := db.RetrieveAccount(decoded[1], values[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
