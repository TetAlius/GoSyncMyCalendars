package frontend

import (
	"github.com/TetAlius/GoSyncMyCalendars/frontend/handlers"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"html/template"
	"net"
	"net/http"
	"strconv"
)

//Frontend object
type Frontend struct {
	IP            net.IP
	Port          int
	googleHandler *handlers.Google
}

//NewFrontend creates a frontend
func NewFrontend(ip string, port int) *Frontend {
	googleHandler := handlers.NewGoogleHandler()
	frontend := Frontend{net.ParseIP(ip), port, googleHandler}
	return &frontend
}

//Start the frontend
func (f *Frontend) Start() error {
	log.Debugln("Start frontend")
	webServerMux := http.NewServeMux()

	cssFileServer := http.StripPrefix("/css/", http.FileServer(http.Dir("./frontend/resources/css/")))
	jsFileServer := http.StripPrefix("/js/", http.FileServer(http.Dir("./frontend/resources/js/")))
	imagesFileServer := http.StripPrefix("/images/", http.FileServer(http.Dir("./frontend/resources/images/")))

	webServerMux.Handle("/css/", cssFileServer)
	webServerMux.Handle("/js/", jsFileServer)
	webServerMux.Handle("/images/", imagesFileServer)
	webServerMux.HandleFunc("/", f.indexHandler)

	webServerMux.HandleFunc("/SignInWithGoogle", f.googleHandler.SignInHandler)

	/*
		http.HandleFunc("/signInWithOutlook", outlookSignInHandler)
		http.HandleFunc("/outlook", outlookTokenHandler)
		http.HandleFunc("/calendars", listCalendarsHandler)
		http.HandleFunc("/SignInWithGoogle", googleSignInHandler)
		http.HandleFunc("/google", googleTokenHandler)
		http.HandleFunc("/signUp", signUpHandler)
		http.HandleFunc("/signIn", signInHandler)
		http.HandleFunc("/cookies", cookiesHandlerTest)
	*/

	laddr := f.IP.String() + ":" + strconv.Itoa(f.Port)
	log.Infof("Web server listening at %s", laddr)

	go func() {
		http.ListenAndServe(laddr, webServerMux)
	}()

	return nil
}

//indexHandler load the index.html web page
//func (s *server) indexHandler(w http.ResponseWriter, r *http.Request, title string) {
func (f *Frontend) indexHandler(w http.ResponseWriter, r *http.Request) {
	// 404 page
	if r.URL.Path != "/" {
		f.errorHandler(w, r, http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("./frontend/welcome.html")
	if err != nil {
		log.Errorln("Error reading config.json: %s", err.Error())
	}

	err = t.Execute(w, nil) //No template at this moment
	if err != nil {
		log.Errorln(err)
	}
}

//errorHandler if something can not be loaded, call the 404 web page
func (f *Frontend) errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	http.ServeFile(w, r, "./frontend/404.html")
}

//Stop the frontend
func (f *Frontend) Stop() error {
	//TODO Complete
	log.Debugln("Stop frontend")
	return nil
}
