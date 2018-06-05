package worker

import (
	"fmt"

	"reflect"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/backend/db"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

const (
	stateInitial = iota
	stateQuit
)

type Worker struct {
	Events   chan api.EventManager
	state    int
	database db.Database
}

func New(maxWorkers int, database db.Database) (worker *Worker) {
	worker = &Worker{Events: make(chan api.EventManager), state: stateInitial, database: database}
	return
}
func (worker *Worker) IsClosed() bool {
	return worker.state == stateQuit
}
func (worker *Worker) Start() {
	worker.Process()
}

func (worker *Worker) Stop() (err error) {
	worker.state = stateQuit
	log.Debugln("closing workers")
	close(worker.Events)
	log.Debugln("close workers")
	return
}

func (worker *Worker) Process() {
	var event api.EventManager
	for worker.state != stateQuit {
		event = <-worker.Events
		worker.processSynchronization(event)

	}
	for event = range worker.Events {
		worker.processSynchronization(event)
	}
}

type SynchronizeError struct {
	State int
	ID    string
}

func (err SynchronizeError) Error() string {
	return fmt.Sprintf("state: %d not suported for event with ID: %s", err.State, err.ID)
}

func (worker *Worker) processSynchronization(event api.EventManager) {
	if event.GetState() == api.Updated && worker.database.EventAlreadyUpdated(event) {
		return
	}
	if event.GetState() == api.Deleted && !worker.database.ExistsEvent(event) {
		return
	}
	switch event.GetState() {
	case api.Created:
		worker.database.SavePrincipalEvent(event)
	case api.Updated:
		worker.database.UpdateModificationDate(event)
	case api.Deleted:
		worker.database.DeleteEvent(event)
	}
	for _, toSync := range event.GetRelations() {
		api.Convert(event, toSync)
		err := worker.synchronizeEvents(event, toSync)
		if err != nil && reflect.TypeOf(err).Kind() != reflect.TypeOf(SynchronizeError{}).Kind() {
			go func() {
				for toSync.CanProcessAgain() {
					toSync.IncrementBackoff()
					err := worker.synchronizeEvents(event, toSync)
					if err != nil {
						continue
					} else {
						//Synchronized correctly
						break
					}
				}
			}()
		} else if err != nil {
			event.MarkWrong()
		}

	}
	return
}

func (worker *Worker) synchronizeEvents(from api.EventManager, to api.EventManager) (err error) {
	switch from.GetState() {
	case api.Created:
		err = worker.createEvent(from, to)
	case api.Updated:
		err = worker.updateEvent(from, to)
	case api.Deleted:
		err = worker.deleteEvent(from, to)
	default:
		return SynchronizeError{State: from.GetState(), ID: from.GetID()}
	}
	return
}

func (worker *Worker) updateEvent(from api.EventManager, to api.EventManager) (err error) {
	err = to.Update()
	if err != nil {
		log.Errorf("error updating event: %s, from event: %s", to.GetID(), from.GetID())
		return err
	}
	return worker.database.UpdateModificationDate(to)

}
func (worker *Worker) createEvent(from api.EventManager, to api.EventManager) (err error) {
	err = to.Create()
	if err != nil {
		log.Errorf("error updating event: %s, from event: %s", to.GetID(), from.GetID())
		return err
	}

	return worker.database.SaveEventsRelation(from, to)

}

func (worker *Worker) deleteEvent(from api.EventManager, to api.EventManager) (err error) {
	err = to.Delete()
	if err != nil {
		log.Errorf("error updating event: %s, from event: %s", to.GetID(), from.GetID())
		return err
	}

	return worker.database.DeleteEvent(to)

}
