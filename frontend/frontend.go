package frontend

import (
	"html/template"
	"net"
	"net/http"
	"strconv"

	"github.com/TetAlius/GoSyncMyCalendars/frontend/handlers"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

//Frontend object
type Server struct {
	IP             net.IP
	Port           int
	googleHandler  *handlers.Google
	outlookHandler *handlers.Outlook
}

//NewFrontend creates a frontend
func NewServer(ip string, port int) *Server {

	googleHandler := handlers.NewGoogleHandler()

	outlookHandler := handlers.NewOutlookHandler()

	server := Server{net.ParseIP(ip), port, googleHandler, outlookHandler}
	return &server
}

//Start the frontend
func (s *Server) Start() error {
	log.Debugln("Start frontend")
	webServerMux := http.NewServeMux()

	cssFileServer := http.StripPrefix("/css/", http.FileServer(http.Dir("./frontend/resources/css/")))
	jsFileServer := http.StripPrefix("/js/", http.FileServer(http.Dir("./frontend/resources/js/")))
	imagesFileServer := http.StripPrefix("/images/", http.FileServer(http.Dir("./frontend/resources/images/")))

	webServerMux.Handle("/css/", cssFileServer)
	webServerMux.Handle("/js/", jsFileServer)
	webServerMux.Handle("/images/", imagesFileServer)
	webServerMux.HandleFunc("/", s.indexHandler)

	webServerMux.HandleFunc("/SignInWithGoogle", s.googleHandler.SignInHandler)

	webServerMux.HandleFunc("/SignInWithOutlook", s.outlookHandler.SignInHandler)
	/*	http.HandleFunc("/calendars", listCalendarsHandler)
		http.HandleFunc("/google", googleTokenHandler)
		http.HandleFunc("/signUp", signUpHandler)
		http.HandleFunc("/signIn", signInHandler)
		http.HandleFunc("/cookies", cookiesHandlerTest)
	*/

	laddr := s.IP.String() + ":" + strconv.Itoa(s.Port)
	log.Infof("Web server listening at %s", laddr)

	go func() {
		http.ListenAndServe(laddr, webServerMux)
	}()

	return nil
}

//indexHandler load the index.html web page
//func (s *server) indexHandler(w http.ResponseWriter, r *http.Request, title string) {
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	// 404 page
	if r.URL.Path != "/" {
		s.errorHandler(w, r, http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("./frontend/resources/html/welcome.html")
	if err != nil {
		log.Errorln("Error reading config.json: %s", err.Error())
	}

	err = t.Execute(w, nil) //No template at this moment
	if err != nil {
		log.Errorln(err)
	}
}

//errorHandler if something can not be loaded, call the 404 web page
func (s *Server) errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	http.ServeFile(w, r, "./frontend/resources/html/404.html")
}

//Stop the frontend
func (s *Server) Stop() error {
	//TODO Complete
	log.Debugln("Stop frontend")
	return nil
}
