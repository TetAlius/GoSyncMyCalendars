package backend

import (
	"net"
	"net/http"
	"strconv"

	"github.com/TetAlius/GoSyncMyCalendars/backend/handlers"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

//Backend object
type Server struct {
	IP             net.IP
	Port           int
	googleHandler  *handlers.Google
	outlookHandler *handlers.Outlook
}

type Accounter interface {
	GetAllCalendars()
	GetPrimaryCalendar()
	GetCalendar(calendarID string)
	CreateCalendar(calendarData []byte)
	UpdateCalendar(calendarID string, calendarData []byte)
	DeleteCalendar(calendarID string)

	GetAllEventsFromCalendar(calendarID string)
	CreateEvent(calendarID string, eventData []byte)
	// Outlook
	//UpdateEvent(eventID string, eventData []byte)
	// Google
	//UpdateEvent(calendarID string, eventID string, eventData []byte)
	//Outlook
	//DeleteEvent(eventID string)
	//Google
	//DeleteEvent(calendarID string, eventID string)
	//Outlook
	// GetEvent(eventID string)
	//Google
	//GetEvent(calendarID string, eventID string)
}

//NewBackend creates a backend
func NewServer(ip string, port int) *Server {
	googleHandler := handlers.NewGoogleHandler()
	outlookHandler := handlers.NewOutlookHandler()
	server := Server{net.ParseIP(ip), port, googleHandler, outlookHandler}
	return &server
}

//Start the backend
func (s *Server) Start() {
	log.Debugln("Start backend")
	webServerMux := http.NewServeMux()

	webServerMux.HandleFunc("/google", s.googleHandler.TokenHandler)
	webServerMux.HandleFunc("/outlook", s.outlookHandler.TokenHandler)

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
