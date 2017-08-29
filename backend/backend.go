package backend

import (
	"net"
	"net/http"
	"strconv"

	"github.com/TetAlius/GoSyncMyCalendars/backend/handlers"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

//Backend object
type Backend struct {
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
func NewBackend(ip string, port int) *Backend {
	googleHandler := handlers.NewGoogleHandler()
	outlookHandler := handlers.NewOutlookHandler()
	backend := Backend{net.ParseIP(ip), port, googleHandler, outlookHandler}
	return &backend
}

//Start the backend
func (b *Backend) Start() {
	log.Debugln("Start backend")
	webServerMux := http.NewServeMux()

	webServerMux.HandleFunc("/google", b.googleHandler.TokenHandler)
	webServerMux.HandleFunc("/outlook", b.outlookHandler.TokenHandler)

	laddr := b.IP.String() + ":" + strconv.Itoa(b.Port)
	log.Infof("Backend server listening at %s", laddr)
	http.HandleFunc("/google", b.googleHandler.TokenHandler)
	err := http.ListenAndServe(laddr, webServerMux)
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}
}

//Stop the backend
func (b *Backend) Stop() error {
	//TODO Complete
	log.Debugln("Stop backend")
	return nil
}
