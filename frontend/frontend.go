package frontend

import (
	"html/template"
	"net"
	"net/http"
	"strings"

	"context"
	"fmt"
	"time"

	"github.com/TetAlius/GoSyncMyCalendars/db"
	"github.com/TetAlius/GoSyncMyCalendars/frontend/handlers"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
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
	User      db.User
	Error     string
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
			log.Errorf("%s", err.Error())
		}
	}()

	return
}

//indexHandler load the index.html web page
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
		serverError(w, err)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Errorln(err)
		serverError(w, err)
		return
	}
}

func (s *Server) calendarListHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := manageSession(w, r)
	if !ok {
		return
	}

	data := PageInfo{
		PageTitle: "Calendars - GoSyncMyCalendars",
		User:      *currentUser,
	}
	t, err := template.ParseFiles("./frontend/resources/html/shared/layout.html", "./frontend/resources/html/calendar-list.html")
	if err != nil {
		log.Errorf("error parsing files: %s", err.Error())
		serverError(w, err)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Errorf("error executing templates: %s", err.Error())
		serverError(w, err)
		return
	}
}

func (s *Server) userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		token := r.FormValue("idtoken")
		email, _, _ := util.MailFromToken(strings.Split(token, "."), "==")
		user, err := db.GetUserFromToken(email) //strings.Split(token, ".")[1])
		if err != nil {
			log.Errorf("error retrieving user: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{Name: "session", Value: user.ID.String()}
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
		serverError(w, err)
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

func serverError(w http.ResponseWriter, error error) {
	t, err := template.ParseFiles("./frontend/resources/html/shared/layout.html", "./frontend/resources/html/500.html")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := PageInfo{
		PageTitle: "Something went wrong :(",
		Error:     error.Error(),
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Errorln(err)
	}

}

func manageSession(w http.ResponseWriter, r *http.Request) (*db.User, bool) {
	session, ok := r.Context().Value("Session").(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil, false
	}
	user, err := db.RetrieveUser(session)
	if err != nil {
		serverError(w, err)
		return nil, false
	}
	if user == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil, false
	}
	return user, true
}

//Stop the frontend
func (s *Server) Stop() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	log.Debugf("Stopping frontend with ctx: %s", ctx)
	err = s.server.Shutdown(nil)
	return
}
