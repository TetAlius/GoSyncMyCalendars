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
	sentry raven.Client
	client *sql.DB
}

func New(client *sql.DB, sentry raven.Client) Database {
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
	//TODO:
	//if err != nil {
	//	log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
	//	return
	//}
	subscriptions = append(subscriptions, subs)
	data.saveSubscription(transaction, subs, calendar)
	events, err := calendar.GetAllEvents()
	//TODO:
	//if err != nil {
	//	log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
	//	return
	//}
	err = data.savePrincipalEvents(transaction, events)
	//TODO:
	//if err != nil {
	//	log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
	//	return
	//}

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
	//TODO:
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

func (data Database) StopSync(principalSubscriptionUUID string, userEmail string, userUUID string) (err error) {
	subscriptions, err := data.RetrieveAllSubscriptionsFromUser(principalSubscriptionUUID, userEmail, userUUID)
	transaction, err := data.client.Begin()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error starting transaction: %s", err.Error())
		return
	}
	for _, subscription := range subscriptions {
		acc := subscription.GetAccount()
		//TODO: manage when account access is refused
		if err = acc.Refresh(); err != nil {
			continue
		}
		go func() { data.UpdateAccountFromUser(acc, userUUID) }()
		//err := subscription.Delete()
		if err != nil {
			log.Errorf("error deleting subscription: %s", err.Error())
		}
		err = data.deleteEventsFromSubscription(transaction, subscription)
		if err != nil {
			log.Errorf("error deleting events from subscription: %s", subscription.GetUUID())
			break
		}
		err = data.deleteSubscription(transaction, subscription)
		if err != nil {
			log.Errorf("error deleting subscription: %s", subscription.GetUUID())
			break
		}
	}
	if err != nil {
		transaction.Rollback()
		log.Errorf("error deleting subscriptions for user: %s", userUUID)
		return
	}
	transaction.Commit()
	return
}
