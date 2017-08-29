package tests

import (
	"github.com/TetAlius/GoSyncMyCalendars/backend"
	"github.com/TetAlius/GoSyncMyCalendars/frontend"
)

func SetupFrontend() *frontend.Server {
	f := frontend.NewServer("127.0.0.1", 8080)
	f.Start()
	return f

}
func SetupBackend() *backend.Server {
	b := backend.NewServer("127.0.0.1", 8081)
	b.Start()
	return b
}
