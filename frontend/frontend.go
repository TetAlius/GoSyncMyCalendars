package frontend

import (
	"html/template"
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
	mux            http.Handler
}

type PageInfo struct {
	PageTitle string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//NewFrontend creates a frontend
func NewServer(ip string, port int) *Server {
	googleHandler := handlers.NewGoogleHandler()
	outlookHandler := handlers.NewOutlookHandler()
	mux := http.NewServeMux()
	server := Server{IP: net.ParseIP(ip), Port: port, googleHandler: googleHandler, outlookHandler: outlookHandler}
	cssFileServer := http.StripPrefix("/css/", http.FileServer(http.Dir("./frontend/resources/css/")))
	jsFileServer := http.StripPrefix("/js/", http.FileServer(http.Dir("./frontend/resources/js/")))
	imagesFileServer := http.StripPrefix("/images/", http.FileServer(http.Dir("./frontend/resources/images/")))

	mux.Handle("/css/", cssFileServer)
	mux.Handle("/js/", jsFileServer)
	mux.Handle("/images/", imagesFileServer)
	mux.HandleFunc("/", server.indexHandler)

	mux.HandleFunc("/SignInWithGoogle", server.googleHandler.SignInHandler)

	mux.HandleFunc("/SignInWithOutlook", server.outlookHandler.SignInHandler)
	mux.HandleFunc("/calendars", server.calendarListHandler)
	mux.HandleFunc("/user", server.userHandler)
	server.mux = AddContext(mux)

	return &server
}
func AddContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("session")
		if cookie != nil {
			ctx := context.WithValue(r.Context(), "Session", cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
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
	data := PageInfo{
		PageTitle: "GoSyncMyCalendars",
	}
	t, err := template.ParseFiles("./frontend/resources/html/shared/layout.html", "./frontend/resources/html/index.html")
	if err != nil {
		log.Errorln("error parsing files %s", err.Error())
		serverError(w)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Errorln(err)
		serverError(w)
		return
	}
}

func (s *Server) calendarListHandler(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("Session")
	user, err := retrieveUser(session)
	if err != nil {
		serverError(w)
		return
	}
	if len(user) == 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := PageInfo{
		PageTitle: "Calendars - GoSyncMyCalendars",
	}
	t, err := template.ParseFiles("./frontend/resources/html/shared/layout.html", "./frontend/resources/html/calendar-list.html")
	if err != nil {
		log.Errorln("error parsing files: %s", err.Error())
		serverError(w)
		return
	}

	err = t.Execute(w, data) //No template at this moment
	if err != nil {
		log.Errorf("error executing templates: %s", err.Error())
		serverError(w)
		return
	}
}

func (s *Server) userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		token := r.FormValue("idtoken")
		log.Debugln(token)
		cookie := http.Cookie{Name: "session", Value: token}
		http.SetCookie(w, &cookie)
		//TODO: handle user
		w.WriteHeader(http.StatusAccepted)
	default:
		notFound(w)
	}
}

func notFound(w http.ResponseWriter) {
	t, err := template.ParseFiles("./frontend/resources/html/shared/layout.html", "./frontend/resources/html/404.html")
	if err != nil {
		serverError(w)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := PageInfo{
		PageTitle: "Not found :(",
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Errorln(err)
	}
}

func serverError(w http.ResponseWriter) {
	t, err := template.ParseFiles("./frontend/resources/html/shared/layout.html", "./frontend/resources/html/500.html")
	if err != nil {
		panic(err) // or do something useful
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := PageInfo{
		PageTitle: "Something went wrong :(",
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Errorln(err)
	}

}

func retrieveUser(session interface{}) (string, error) {
	//TODO:
	return "asd", nil

}

//Stop the frontend
func (s *Server) Stop() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	log.Debugf("Stopping frontend with ctx: %s", ctx)
	err = s.server.Shutdown(nil)
	return
}
