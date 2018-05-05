package backend

import (
	"net"
	"net/http"

	"fmt"

	"context"
	"time"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

//Backend object
type Server struct {
	IP      net.IP
	Port    int
	workers *[]Worker
	server  *http.Server
	mux     *http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//NewBackend creates a backend
func NewServer(ip string, port int) *Server {
	workers := new([]Worker)
	server := Server{IP: net.ParseIP(ip), Port: port, workers: workers, mux: http.NewServeMux()}
	server.mux.HandleFunc("/google", server.GoogleTokenHandler)
	server.mux.HandleFunc("/google/watcher", server.GoogleWatcherHandler)
	server.mux.HandleFunc("/outlook", server.OutlookTokenHandler)
	server.mux.HandleFunc("/outlook/watcher", server.OutlookWatcherHandler)
	return &server
}

//Start the backend
func (s *Server) Start() (err error) {
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
	return
}
