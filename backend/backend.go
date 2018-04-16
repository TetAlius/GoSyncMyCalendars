package backend

import (
	"net"
	"net/http"
	"strconv"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

//Backend object
type Server struct {
	IP      net.IP
	Port    int
	workers *[]Worker
}

//NewBackend creates a backend
func NewServer(ip string, port int) *Server {
	workers := new([]Worker)
	server := Server{net.ParseIP(ip), port, workers}
	return &server
}

//Start the backend
func (s *Server) Start() {
	log.Debugln("Start backend")
	webServerMux := http.NewServeMux()

	webServerMux.HandleFunc("/google", s.GoogleTokenHandler)
	webServerMux.HandleFunc("/google/watcher", s.GoogleWatcherHandler)
	webServerMux.HandleFunc("/outlook", s.OutlookTokenHandler)
	webServerMux.HandleFunc("/outlook/watcher", s.OutlookWatcherHandler)

	laddr := s.IP.String() + ":" + strconv.Itoa(s.Port)
	log.Infof("Backend server listening at %s", laddr)

	err := http.ListenAndServe(laddr, webServerMux)
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}

}

//Stop the backend
func (s *Server) Stop() error {
	//TODO Complete
	log.Debugln("TODO: Stop backend")
	return nil
}
