package db

import (
	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/getsentry/raven-go"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Database struct {
	sentry *raven.Client
	client *sql.DB
}

func New(client *sql.DB, sentry *raven.Client) Database {
	return Database{
		client: client,
		sentry: sentry,
	}

}

func (data Database) Close() error {
	return data.client.Close()
}

//TODO: raven this
func (data Database) StartSync(calendar api.CalendarManager, userUUID string) (err error) {
	transaction, err := data.client.Begin()
	if err != nil {
		log.Errorf("error creating transaction: %s", err.Error())
		return
	}
	data.UpdateAccountFromUser(calendar.GetAccount(), userUUID)
	data.UpdateCalendarFromUser(calendar, userUUID)
	var subscriptions []api.SubscriptionManager
	var subs api.SubscriptionManager
	switch calendar.(type) {
	case *api.GoogleCalendar:
		subs = api.NewGoogleSubscription(uuid.New().String())
		err = subs.Subscribe(calendar)
	case *api.OutlookCalendar:
		subs = api.NewOutlookSubscription(uuid.New().String())
		err = subs.Subscribe(calendar)
	}
	//if err != nil {
	//	log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
	//	return
	//}
	subscriptions = append(subscriptions, subs)
	data.saveSubscription(transaction, subs, calendar)
	events, err := calendar.GetAllEvents()
	data.savePrincipalEvents(transaction, events)

	var eventsCreated []api.EventManager
	for _, cal := range calendar.GetCalendars() {
		data.UpdateAccountFromUser(cal.GetAccount(), userUUID)
		data.UpdateCalendarFromUser(cal, userUUID)
		var subscript api.SubscriptionManager
		var toEvent api.EventManager
		switch cal.(type) {
		case *api.GoogleCalendar:
			toEvent = &api.GoogleEvent{}
			subscript = api.NewGoogleSubscription(uuid.New().String())
		case *api.OutlookCalendar:
			toEvent = &api.OutlookEvent{}
			subscript = api.NewOutlookSubscription(uuid.New().String())
		}
		for _, event := range events {
			api.Convert(event, toEvent)
			err = toEvent.SetCalendar(cal)
			if err != nil {
				log.Errorf("error converting event for calendar: %s, error: %s", cal.GetUUID(), err.Error())
				break
			}
			toEvent.Create()
			if err != nil {
				log.Errorf("error creating event for calendar: %s, error: %s", cal.GetUUID(), err.Error())
				break
			}
			eventsCreated = append(eventsCreated, toEvent)
			err = data.saveEventsRelation(transaction, event, toEvent)
			if err != nil {

			}
		}
		if err != nil {
			log.Errorln("some error saving events")
			break
		}
		subscript.Subscribe(cal)
		if err != nil {
			log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
			break
		}
		subscriptions = append(subscriptions, subscript)
		data.saveSubscription(transaction, subscript, cal)
	}
	//if err != nil {
	//	transaction.Rollback()
	//	for _, subscription := range subscriptions {
	//		subscription.Delete()
	//	}
	//
	//	for _, event := range eventsCreated {
	//		event.Delete()
	//	}
	//	return
	//}
	transaction.Commit()
	return
}
