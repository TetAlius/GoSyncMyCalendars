package worker

import (
	"fmt"

	"reflect"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/backend/db"
	"github.com/TetAlius/GoSyncMyCalendars/convert"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

const (
	stateInitial = iota
	stateQuit
)

// Object that manages the different kinds of synchronization
type Worker struct {
	Events   chan api.EventManager
	state    int
	database db.Database
}

// Function that returns a new worker from given info
func New(maxWorkers int, database db.Database) (worker *Worker) {
	worker = &Worker{Events: make(chan api.EventManager), state: stateInitial, database: database}
	return
}

// Method that returns whether the channel is closed
func (worker *Worker) IsClosed() bool {
	return worker.state == stateQuit
}

// Method that starts the processing of the worker
func (worker *Worker) Start() {
	worker.Process()
}

// Method that stops the processing of the worker
func (worker *Worker) Stop() (err error) {
	worker.state = stateQuit
	log.Debugln("closing workers")
	close(worker.Events)
	log.Debugln("close workers")
	return
}

// Method that process all requests of sync
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

// Specific synchronization error
type SynchronizeError struct {
	State int
	ID    string
}

// Method to implement the interface error
func (err SynchronizeError) Error() string {
	return fmt.Sprintf("state: %d not suported for event with ID: %s", err.State, err.ID)
}

// Method that process a specific request of sync
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
		err := toSync.GetCalendar().GetAccount().Refresh()
		go worker.database.UpdateAccount(toSync.GetCalendar().GetAccount())
		err = worker.synchronizeEvents(event, toSync)
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
		}

	}
	return
}

// Method that synchronize to events. If the request gets here, all database checks have passed
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

// Method that manages an update
func (worker *Worker) updateEvent(from api.EventManager, to api.EventManager) (err error) {
	convert.Convert(from, to)
	err = to.Update()
	if err != nil {
		log.Errorf("error updating event: %s, from event: %s", to.GetID(), from.GetID())
		return err
	}
	return worker.database.UpdateModificationDate(to)

}

// Method that manages a creation
func (worker *Worker) createEvent(from api.EventManager, to api.EventManager) (err error) {
	convert.Convert(from, to)
	err = to.Create()
	if err != nil {
		log.Errorf("error updating event: %s, from event: %s", to.GetID(), from.GetID())
		return err
	}

	return worker.database.SaveEventsRelation(from, to)

}

// Method that manages a deletion
func (worker *Worker) deleteEvent(from api.EventManager, to api.EventManager) (err error) {
	if !worker.database.ExistsEvent(to) {
		return nil
	}
	err = to.Delete()
	if err != nil {
		log.Errorf("error updating event: %s, from event: %s", to.GetID(), from.GetID())
		return err
	}

	return worker.database.DeleteEvent(to)

}
