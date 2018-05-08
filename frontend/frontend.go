package frontend

import (
	"html/template"
	"io/ioutil"
	"net"
	"net/http"

	"context"
	"fmt"
	"time"

	"github.com/TetAlius/GoSyncMyCalendars/frontend/handlers"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

//Frontend object
type Server struct {
	IP             net.IP
	Port           int
	googleHandler  *handlers.Google
	outlookHandler *handlers.Outlook
	server         *http.Server
	mux            *http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//NewFrontend creates a frontend
func NewServer(ip string, port int) *Server {
	googleHandler := handlers.NewGoogleHandler()
	outlookHandler := handlers.NewOutlookHandler()
	server := Server{IP: net.ParseIP(ip), Port: port, googleHandler: googleHandler, outlookHandler: outlookHandler, mux: http.NewServeMux()}
	cssFileServer := http.StripPrefix("/css/", http.FileServer(http.Dir("./frontend/resources/css/")))
	jsFileServer := http.StripPrefix("/js/", http.FileServer(http.Dir("./frontend/resources/js/")))
	imagesFileServer := http.StripPrefix("/images/", http.FileServer(http.Dir("./frontend/resources/images/")))

	server.mux.Handle("/css/", cssFileServer)
	server.mux.Handle("/js/", jsFileServer)
	server.mux.Handle("/images/", imagesFileServer)
	server.mux.HandleFunc("/", server.indexHandler)

	server.mux.HandleFunc("/SignInWithGoogle", server.googleHandler.SignInHandler)

	server.mux.HandleFunc("/SignInWithOutlook", server.outlookHandler.SignInHandler)

	return &server
}

//Start the frontend
func (s *Server) Start() (err error) {
	log.Debugln("Start frontend")

	laddr := fmt.Sprintf("%s:%d", s.IP.String(), s.Port)
	h := &http.Server{Addr: fmt.Sprintf(":%d", s.Port), Handler: s}
	s.server = h
	go func() {
		log.Infof("Web server listening at %s", laddr)

		if err := s.server.ListenAndServe(); err != nil {
			log.Fatalf("%s", err.Error())

		}
	}()

	return
}

//indexHandler load the index.html web page
//func (s *server) indexHandler(w http.ResponseWriter, r *http.Request, title string) {
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	// 404 page
	if r.URL.Path != "/" {
		notFound(w)
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
		panic(err) // or do something useful
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(contents)

}

//Stop the frontend
func (s *Server) Stop() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	log.Debugf("Stopping frontend with ctx: %s", ctx)
	err = s.server.Shutdown(nil)
	return
}
