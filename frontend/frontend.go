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

	"database/sql"

	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"

	"os"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	"github.com/TetAlius/GoSyncMyCalendars/frontend/db"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/getsentry/raven-go"
	"github.com/google/uuid"
)

//Frontend object
type Server struct {
	IP       net.IP
	Port     int
	server   *http.Server
	mux      http.Handler
	database db.Database
	sentry   *raven.Client
}

type PageInfo struct {
	PageTitle string
	User      db.User
	Account   db.Account
	Calendars []db.Calendar
	Error     string
}

var root string

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	}, "dec": func(i int) int {
		return i - 1
	}, "existsUUID": func(uid uuid.UUID) bool {
		return uid.String() != "00000000-0000-0000-0000-000000000000"
	}, "endpoint": func() string {
		return os.Getenv("ENDPOINT")
	}, "disabledByUUID": func(uid uuid.UUID) string {
		if uid.String() != "00000000-0000-0000-0000-000000000000" {
			return "disabled"
		}
		return ""
	},
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//NewFrontend creates a frontend
func NewServer(ip string, port int, dir string, database *sql.DB, sentry *raven.Client) *Server {
	mux := http.NewServeMux()
	root = dir
	server := Server{IP: net.ParseIP(ip), Port: port, database: db.New(database, sentry), sentry: sentry}
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

		if err := s.server.ListenAndServeTLS("server.crt", "server.key"); err != nil && err != http.ErrServerClosed {
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
	t, err := template.New("layout.html").Funcs(funcMap).ParseFiles(root+"/html/shared/layout.html", root+"/html/index.html")
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
	currentUser, ok := s.manageSession(w, r)
	if !ok {
		return
	}
	err := s.database.SetUserAccounts(currentUser)
	if err != nil {
		serverError(w, err)
		return
	}

	data := PageInfo{
		PageTitle: "Calendars",
		User:      *currentUser,
		Account:   currentUser.PrincipalAccount,
	}
	t, err := template.New("layout.html").Funcs(funcMap).ParseFiles(root+"/html/shared/layout.html", root+"/html/calendars/list.html")
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
	currentUser, ok := s.manageSession(w, r)
	if !ok {
		return
	}
	err := s.database.SetUserAccounts(currentUser)
	if err != nil {
		serverError(w, err)
		return
	}

	data := PageInfo{
		PageTitle: "Accounts",
		User:      *currentUser,
	}
	t, err := template.New("layout.html").Funcs(funcMap).ParseFiles(root+"/html/shared/layout.html", root+"/html/accounts/list.html")
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
	currentUser, ok := s.manageSession(w, r)
	if !ok {
		return
	}
	param := r.URL.Path[len("/accounts/"):]
	id, err := strconv.Atoi(param)
	if err != nil {
		serverError(w, err)
		return
	}
	account, err := s.database.FindAccount(currentUser, id)
	if err != nil {
		serverError(w, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		// Form submitted
		r.ParseForm() // Required if you don't call r.FormValue()
		calendarIDs := r.Form["calendars"]
		log.Debugf("IDs: %s", calendarIDs)
		s.database.AddCalendarsToAccount(currentUser, account, calendarIDs)
	default:
		serverError(w, err)
		return
	}
	err = s.database.FindCalendars(&account)
	if err != nil {
		serverError(w, err)
		return
	}

	data := PageInfo{
		PageTitle: account.Email,
		User:      *currentUser,
		Account:   account,
	}

	t, err := template.New("layout.html").Funcs(funcMap).ParseFiles(root+"/html/shared/layout.html", root+"/html/accounts/show.html")
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
	currentUser, ok := s.manageSession(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	id := r.URL.Path[len("/calendars/"):]
	switch r.Method {
	case http.MethodDelete:

		err := s.database.DeleteCalendar(currentUser, id)
		if err != nil {
			serverError(w, err)
			return
		}
	case http.MethodPost:
		log.Debugf("id: %s", id)
		r.ParseForm() // Required if you don't call r.FormValue()
		calendarIDs := r.Form["calendars"]
		log.Debugf("calendar ids: %s", calendarIDs)

		err := s.database.AddCalendarsRelation(currentUser, id, calendarIDs)
		if err != nil {
			serverError(w, err)
			return
		}
		http.Redirect(w, r, "/calendars", http.StatusFound)
	case http.MethodPatch:
		parent := r.FormValue("parent")
		log.Debugf("parents ids: %s", parent)

		err := s.database.UpdateCalendar(currentUser, id, parent)
		if err != nil {
			serverError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		notFound(w)
		return
	}

}

func (s *Server) userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		token := r.FormValue("idtoken")
		user, err := userFromToken(strings.Split(token, "."))
		err = s.database.FindOrCreateUser(user)
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
	t, err := template.New("layout.html").Funcs(funcMap).ParseFiles(root+"/html/shared/layout.html", root+"/html/404.html")
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
	t, err := template.New("layout.html").Funcs(funcMap).ParseFiles(root+"/html/shared/layout.html", root+"/html/500.html")
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

func (s *Server) manageSession(w http.ResponseWriter, r *http.Request) (*db.User, bool) {
	session, ok := r.Context().Value("Session").(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil, false
	}
	user, err := s.database.RetrieveUser(session)
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
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	log.Debugf("Stopping frontend with ctx: %s", ctx)
	var returnErr error
	err = s.server.Shutdown(ctx)
	if err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"stopping": "frontend server"})
		returnErr = err
	}
	err = s.database.Close()
	if err != nil {
		raven.CaptureErrorAndWait(err, map[string]string{"stopping": "frontend database"})
		returnErr = err
	}
	return returnErr
}
func userFromToken(tokens []string) (user *db.User, err error) {
	if len(tokens) < 2 {
		return nil, errors.New("TokenID was not parsed correctly")
	}
	encodedToken := strings.Replace(
		strings.Replace(tokens[1], "-", "_", -1),
		"+", "/", -1)
	if encodedToken[len(encodedToken)-1:] != "=" {
		encodedToken = encodedToken + "="
	}
	if encodedToken[len(encodedToken)-2:] != "==" {
		encodedToken = encodedToken + "="
	}
	decodedToken, err := base64.StdEncoding.DecodeString(encodedToken)
	log.Debugf("%s", decodedToken)

	if err != nil {
		log.Errorf("Error decoding token: %s", err.Error())
	}
	var f interface{}
	err = json.Unmarshal(decodedToken, &f)
	if err != nil {
		log.Errorf("Error unmarshaling decoded token: %s", err.Error())
	}

	var email string
	var name string
	if reflect.TypeOf(f) != nil && reflect.TypeOf(f).Kind() == reflect.Map {
		m := f.(map[string]interface{})

		if m["email"] != nil {
			email = m["email"].(string)
		}
		if m["name"] != nil {
			name = m["name"].(string)
		}
		user = &db.User{Email: email, Name: name}
	} else {
		err = customErrors.DecodedError{Message: "Decoded token is not a map"}
	}
	return
}
