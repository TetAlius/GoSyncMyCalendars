package frontend

import (
	"html/template"
	"net"
	"net/http"
	"strings"

	"context"
	"fmt"
	"time"

	"strconv"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/db"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/TetAlius/GoSyncMyCalendars/util"
)

//Frontend object
type Server struct {
	IP     net.IP
	Port   int
	server *http.Server
	mux    http.Handler
}

type PageInfo struct {
	PageTitle string
	User      db.User
	Error     string
	Account   api.AccountManager
	Calendar  api.CalendarManager
	Calendars []api.CalendarManager
}

var root string

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//NewFrontend creates a frontend
func NewServer(ip string, port int, dir string) *Server {
	mux := http.NewServeMux()
	root = dir
	server := Server{IP: net.ParseIP(ip), Port: port}
	cssFileServer := http.StripPrefix("/css/", http.FileServer(http.Dir(root+"/css/")))
	jsFileServer := http.StripPrefix("/js/", http.FileServer(http.Dir(root+"/js/")))
	imagesFileServer := http.StripPrefix("/images/", http.FileServer(http.Dir(root+"/images/")))

	mux.Handle("/css/", cssFileServer)
	mux.Handle("/js/", jsFileServer)
	mux.Handle("/images/", imagesFileServer)
	mux.HandleFunc("/", server.indexHandler)

	mux.HandleFunc("/SignInWithGoogle", server.googleSignInHandler)
	mux.HandleFunc("/google", server.googleTokenHandler)

	mux.HandleFunc("/SignInWithOutlook", server.outlookSignInHandler)
	mux.HandleFunc("/outlook", server.OutlookTokenHandler)

	mux.HandleFunc("/calendars", server.calendarListHandler)
	mux.HandleFunc("/calendars/", server.calendarHandler)
	mux.HandleFunc("/accounts", server.accountListHandler)
	mux.HandleFunc("/accounts/", server.accountHandler)
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
	t, err := template.ParseFiles(root+"/html/shared/layout.html", root+"/html/index.html")
	if err != nil {
		log.Errorln("error parsing files %s", err.Error())
		serverError(w, err)
		return
	}

	err = t.Execute(w, PageInfo{})
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
		PageTitle: "Calendars",
		User:      *currentUser,
	}
	t, err := template.ParseFiles(root+"/html/shared/layout.html", root+"/html/calendar-list.html")
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

func (s *Server) accountListHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := manageSession(w, r)
	if !ok {
		return
	}

	data := PageInfo{
		PageTitle: "Accounts",
		User:      *currentUser,
	}
	t, err := template.ParseFiles(root+"/html/shared/layout.html", root+"/html/accounts/list.html")
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

func (s *Server) accountHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := manageSession(w, r)
	if !ok {
		return
	}
	param := r.URL.Path[len("/accounts/"):]
	id, err := strconv.Atoi(param)
	if err != nil {
		serverError(w, err)
		return
	}
	account, err := currentUser.FindAccount(id)
	if err != nil {
		serverError(w, err)
		return
	}
	if account == nil {
		notFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		// Form submitted
		r.ParseForm() // Required if you don't call r.FormValue()
		calendarIDs := r.Form["calendars"]
		log.Debugf("IDs: %s", calendarIDs)
		currentUser.AddCalendarsToAccount(account, calendarIDs)
	default:
		serverError(w, err)
		return
	}
	account.Refresh()
	calendars, err := currentUser.RetrieveCalendarsFromAccount(account)
	if err != nil {
		serverError(w, err)
		return
	}
	data := PageInfo{
		PageTitle: account.Mail(),
		Account:   account,
		Calendars: calendars,
	}

	t, err := template.ParseFiles(root+"/html/shared/layout.html", root+"/html/accounts/show.html")
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

func (s *Server) calendarHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := manageSession(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		id := r.URL.Path[len("/calendars/"):]
		err := currentUser.DeleteCalendar(id)
		if err != nil {
			serverError(w, err)
			return
		}
	default:
		notFound(w)
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
		cookie := http.Cookie{Name: "session", Value: user.UUID.String()}
		http.SetCookie(w, &cookie)
		//TODO: handle user
		w.WriteHeader(http.StatusAccepted)
	default:
		notFound(w)
	}
}

func notFound(w http.ResponseWriter) {
	t, err := template.ParseFiles(root+"/html/shared/layout.html", root+"/html/404.html")
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
	t, err := template.ParseFiles(root+"/html/shared/layout.html", root+"/html/500.html")
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
