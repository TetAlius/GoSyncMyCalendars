package db

import (
	"crypto/rand"
	"database/sql"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/convert"
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

func randToken() string {
	b := make([]byte, 10)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

//TODO: raven this
func (data Database) StartSync(calendar api.CalendarManager, userUUID string) (err error) {
	var subscriptions []api.SubscriptionManager
	var subs api.SubscriptionManager
	var eventsCreated []api.EventManager
	var events []api.EventManager
	transaction, err := data.client.Begin()
	if err != nil {
		log.Errorf("error creating transaction: %s", err.Error())
		return
	}
	data.UpdateAccountFromUser(calendar.GetAccount(), userUUID)
	events, err = calendar.GetAllEvents()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
		goto End
	}
	data.UpdateCalendarFromUser(calendar, userUUID)
	switch calendar.(type) {
	case *api.GoogleCalendar:
		subs = api.NewGoogleSubscription(uuid.New().String())
	case *api.OutlookCalendar:
		subs = api.NewOutlookSubscription()
	}
	err = subs.Subscribe(calendar)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
		goto End
	}
	subscriptions = append(subscriptions, subs)
	data.saveSubscription(transaction, subs, calendar)
	err = data.savePrincipalEvents(transaction, events)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
		goto End
	}

	for _, cal := range calendar.GetCalendars() {
		data.UpdateAccountFromUser(cal.GetAccount(), userUUID)
		data.UpdateCalendarFromUser(cal, userUUID)
		var subscript api.SubscriptionManager
		switch cal.(type) {
		case *api.GoogleCalendar:
			subscript = api.NewGoogleSubscription(uuid.New().String())
		case *api.OutlookCalendar:
			subscript = api.NewOutlookSubscription()
		}
		for _, event := range events {
			var toEvent api.EventManager
			switch cal.(type) {
			case *api.GoogleCalendar:
				toEvent = &api.GoogleEvent{}
			case *api.OutlookCalendar:
				toEvent = &api.OutlookEvent{}
			}
			convert.Convert(event, toEvent)
			err = toEvent.SetCalendar(cal)
			if err != nil {
				data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
				log.Errorf("error converting event for calendar: %s, error: %s", cal.GetUUID(), err.Error())
				goto End
			}
			err = toEvent.Create()
			if err != nil {
				data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
				log.Errorf("error creating event for calendar: %s, error: %s", cal.GetUUID(), err.Error())
				goto End
			}
			eventsCreated = append(eventsCreated, toEvent)
			err = data.saveEventsRelation(transaction, event, toEvent)
			if err != nil {
				data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
				log.Errorf("error saving relation on database: %s, error: %s", event.GetID(), err.Error())
				goto End
			}
		}
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			log.Errorln("some error saving events")
			goto End
		}
		err := subscript.Subscribe(cal)
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
			goto End
		}
		subscriptions = append(subscriptions, subscript)
		err = data.saveSubscription(transaction, subscript, cal)
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			log.Errorf("error saving subscription to db: %s", subscript.GetID())
			goto End
		}
	}
End:
	if err != nil {
		transaction.Rollback()
		for _, subscription := range subscriptions {
			subscription.Delete()
		}

		for _, event := range eventsCreated {
			event.Delete()
		}
		return
	}
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
		err = subscription.Delete()
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
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
