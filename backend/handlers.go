package main

import "net/http"

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/welcome.html")
}
