package worker

import (
	"fmt"

	"reflect"

	"github.com/TetAlius/GoSyncMyCalendars/api"
)

const (
	stateInitial = iota
	stateMain
	stateExit
	stateLaunch
	stateProcessing
	stateFinish
	stateSleep
	stateWait
	stateTimeout
	stateQuit
)

type Worker struct {
	Events chan api.EventManager
	state  int
}

func New(maxWorkers int) (worker *Worker) {
	worker = &Worker{Events: make(chan api.EventManager, maxWorkers), state: stateInitial}
	return
}
func (worker *Worker) IsClosed() bool {
	return worker.state == stateQuit
}
func (worker *Worker) Start() {
	go worker.Process()
}

func (worker *Worker) Stop() (err error) {
	worker.state = stateQuit
	close(worker.Events)
	return
}

func (worker *Worker) Process() {
	var event api.EventManager
	for worker.state != stateQuit {
		event = <-worker.Events
		processSynchronization(event)

	}
	for event = range worker.Events {
		processSynchronization(event)
	}
}

type SynchronizeError struct {
	State int
	ID    string
}

func (err SynchronizeError) Error() string {
	return fmt.Sprintf("state: %d not suported for event with ID: %s", err.State, err.ID)
}

func processSynchronization(event api.EventManager) {
	for _, toSync := range event.GetRelations() {
		api.Convert(event, toSync)
		err := synchronizeEvents(event, toSync)
		if err != nil && reflect.TypeOf(err).Kind() != reflect.TypeOf(SynchronizeError{}).Kind() {
			go func() {
				for toSync.CanProcessAgain() {
					toSync.IncrementBackoff()
					synchronizeEvents(event, toSync)
				}
			}()
		} else if err != nil {
			event.MarkWrong()
		}
	}
	return
}

func synchronizeEvents(from api.EventManager, to api.EventManager) (err error) {
	switch from.GetState() {
	case api.Created:
		err = to.Create()
	case api.Updated:
		err = to.Update()
	case api.Deleted:
		err = to.Delete()
	default:
		return SynchronizeError{State: from.GetState(), ID: from.GetID()}
	}
	return
}
