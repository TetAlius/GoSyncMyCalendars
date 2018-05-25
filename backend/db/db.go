package db

import (
	"database/sql"
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func (data Database) StartSync(calendar api.CalendarManager, userUUID string) (err error) {
	data.UpdateAccountFromUser(calendar.GetAccount(), userUUID)
	data.UpdateCalendarFromUser(calendar, userUUID)
	var subscriptions []api.SubscriptionManager
	var subs api.SubscriptionManager
	switch calendar.(type) {
	case *api.GoogleCalendar:
		//TODO: change this IDS
		subs = api.NewGoogleSubscription(uuid.New().String(), "URL")
		err = subs.Subscribe(calendar)
	case *api.OutlookCalendar:
		subs = api.NewOutlookSubscription(uuid.New().String(), "URL", "Created,Deleted,Updated")
		err = subs.Subscribe(calendar)
	}
	if err != nil {
		log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
		return
	}
	subscriptions = append(subscriptions, subs)
	data.SaveSubscription(subs, calendar)
	for _, cal := range calendar.GetCalendars() {
		data.UpdateAccountFromUser(cal.GetAccount(), userUUID)
		data.UpdateCalendarFromUser(cal, userUUID)
		var subscript api.SubscriptionManager
		switch cal.(type) {
		case *api.GoogleCalendar:
			//TODO: change this IDs
			subscript = api.NewGoogleSubscription(uuid.New().String(), "URL")
			err = subscript.Subscribe(cal)
		case *api.OutlookCalendar:
			subscript = api.NewOutlookSubscription(uuid.New().String(), "URL", "Created,Deleted,Updated")
			err = subscript.Subscribe(cal)
		}
		if err != nil {
			log.Errorf("error creating subscription for calendar: %s, error: %s", calendar.GetUUID(), err.Error())
			break
		}
		subscriptions = append(subscriptions, subscript)
		data.SaveSubscription(subscript, cal)
	}
	if err != nil {
		for _, subscription := range subscriptions {
			subscription.Delete()
			data.DeleteSubscription(subscription)
		}
	}
	return
}

type WrongKindError struct {
	Mail string
}

func (err *WrongKindError) Error() string {
	return fmt.Sprintf("wrong kind of account %s", err.Mail)
}
