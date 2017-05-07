package backend

import (
	"net"
	"net/http"
	"strconv"

	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

//Backend object
type Backend struct {
	IP   net.IP
	Port int
}

//NewBackend creates a backend
func NewBackend(ip string, port int) *Backend {
	backend := Backend{net.ParseIP(ip), port}
	return &backend
}

//Start the backend
func (b *Backend) Start() {
	log.Debugln("Start backend")
	laddr := b.IP.String() + ":" + strconv.Itoa(b.Port)
	log.Infof("Backend server listening at %s", laddr)
	err := http.ListenAndServe(laddr, nil)
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
