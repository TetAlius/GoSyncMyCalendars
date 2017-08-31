package main

type calendarInfo struct {
	account string
	names   []string
}

/*
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/welcome.html")
}

func listCalendarsHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./frontend/calendars.html")
	if err != nil {
		log.Fatalf("Parse file error: %s", err.Error())
	}

	calendars := []calendarInfo{
		{"outlook@outlook.com", []string{"a", "b"}},
		{"outlook@outlook.com", []string{"a", "b"}},
	}
	log.Debugln(calendars)

	t.Execute(w, calendars)
}
func signUpHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugln(r.Method)
	switch r.Method {
	case "GET":
		log.Debugln("Serving form to sign up")
		http.ServeFile(w, r, "./frontend/sign-up.html")
		break
	case "POST":
		log.Debugln("Storing user info")

		user.CheckInfo(r.FormValue("name"), r.FormValue("surname"), r.FormValue("email"), r.FormValue("pswd1"), r.FormValue("pswd2"))
		log.Debugf("Name: %s, Surname: %s, Mail: %s", r.FormValue("name"), r.FormValue("surname"), r.FormValue("email"))

		u := user.User{Name: r.FormValue("name"), Surname: r.FormValue("surname"), Email: r.FormValue("email"), Pswd: r.FormValue("pswd1")}
		u.Save()
		//user.Save(r.FormValue("name"), r.FormValue("surname"), r.FormValue("email"), r.FormValue("pswd1"))
		http.ServeFile(w, r, "./frontend/welcome.html")
	}
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugln(r.Method)
	switch r.Method {
	case "GET":
		log.Debugln("Serving form to log in")
		http.ServeFile(w, r, "./frontend/sign-in.html")
		break
	case "POST":
		log.Debugln("Cheking user info")
		exists, err := user.CheckExistingUser(r.FormValue("email"), r.FormValue("password"))
		if err != nil {
			log.Errorln("Something went wrong")
		}
		if !exists {
			//TODO
			log.Errorln("User not found on database")
		}
		c, err := r.Cookie("session")
		if err != nil {
			log.Debugln("Cookie was not present. Creating new one.")

			cookie := http.Cookie{Name: "session", Value: "test", Expires: time.Now().Add(24 * time.Hour)}
			http.SetCookie(w, &cookie)
		} else {
			c.Expires = time.Now().Add(24 * time.Hour)
		}

		break
	}
}

func cookiesHandlerTest(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("test")
	log.Debugln(c, err)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "test", Value: "abcd", Expires: expiration}
	http.SetCookie(w, &cookie)
}
*/
