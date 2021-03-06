package backend

import (
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func (s *Server) manageSynchronizationOutlook(notifications []api.OutlookSubscriptionNotification) (err error) {
	tags := map[string]string{"sync": "outlook"}
	for _, subscription := range notifications {
		calendar, err := s.retrieveCalendar(subscription.SubscriptionID, tags)
		if err != nil {
			return err
		}
		if calendar == nil && err == nil {
			return nil
		}
		calendar.GetAccount().Refresh()
		go s.database.UpdateAccount(calendar.GetAccount())
		tags["event"] = subscription.ChangeType
		if subscription.ChangeType == "Missed" {
			err = s.manageByCalendar(calendar, subscription.SubscriptionID, tags)
			if err != nil {
				return err
			}
			continue
		}
		err = s.manageSubscription(calendar, subscription.SubscriptionID, subscription.Data.ID, tags)
		if err != nil {
			log.Errorf("error managing outlook subscription ID: %s", subscription.SubscriptionID)
			return err
		}
	}
	return err
}

func (s *Server) manageByCalendar(calendar api.CalendarManager, subscriptionID string, tags map[string]string) (err error) {
	eventIDs := make(map[string]string)
	events, err := calendar.GetAllEvents()
	if err != nil {
		log.Errorf("error getting all events from cloud: %s", err.Error())
		s.sentry.CaptureErrorAndWait(err, tags)
		return err
	}
	for _, event := range events {
		eventIDs[event.GetID()] = event.GetID()
	}
	IDs, err := s.database.GetEventIDs(subscriptionID)
	if err != nil {
		log.Errorf("error getting all events from db: %s", err.Error())
		s.sentry.CaptureErrorAndWait(err, tags)
		return err
	}
	for _, eventID := range IDs {
		eventIDs[eventID] = eventID
	}
	for eventID := range eventIDs {
		err = s.manageSubscription(calendar, subscriptionID, eventID, tags)
		if err != nil {
			log.Errorf("error managing subscription ID: %s", subscriptionID)
			return err
		}
	}
	return
}

func (s *Server) manageSynchronizationGoogle(subscriptionID string) (err error) {
	tags := map[string]string{"sync": "google"}
	calendar, err := s.retrieveCalendar(subscriptionID, tags)
	if err != nil {
		return err
	}
	if calendar == nil && err == nil {
		return nil
	}
	calendar.GetAccount().Refresh()
	go s.database.UpdateAccount(calendar.GetAccount())
	return s.manageByCalendar(calendar, subscriptionID, tags)
}

func (s *Server) manageSubscription(calendar api.CalendarManager, subscriptionID string, eventID string, tags map[string]string) (err error) {
	onCloud := true
	event, err := calendar.GetEvent(eventID)
	if _, ok := err.(*customErrors.NotFoundError); ok {
		onCloud = false
		err = nil
		event = calendar.CreateEmptyEvent(eventID)
	}
	if err != nil {
		s.sentry.CaptureErrorAndWait(err, tags)
		log.Errorf("error retrieving event from account: %s", err.Error())
		return err
	}
	events, onDB, err := s.database.RetrieveSyncedEventsWithSubscription(eventID, subscriptionID, calendar)
	if err != nil {
		s.sentry.CaptureErrorAndWait(err, tags)
		log.Errorf("error retrieving events synced: %s", err.Error())
		return err
	}
	if !onCloud && !onDB {
		log.Warningf("event with id: %s already deleted", eventID)
		return nil
	}
	if onCloud && onDB && s.database.EventAlreadyUpdated(event) {
		return nil
	}

	event.SetRelations(events)
	state := api.GetChangeType(onCloud, onDB)
	if state == 0 {
		err = fmt.Errorf("synchronization not supported for event: %s", eventID)
		s.sentry.CaptureErrorAndWait(err, tags)
		return err
	}

	event.SetState(state)
	s.worker.Events <- event
	return
}

func (s *Server) retrieveCalendar(subscriptionID string, tags map[string]string) (calendar api.CalendarManager, err error) {
	ok, err := s.database.ExistsSubscriptionFromID(subscriptionID)
	if err != nil && ok {
		//Ignore this subscription
		//Perhaps something went wrong when deleting subscription...
		//s.sentry.CaptureMessageAndWait(fmt.Sprintf("outlook subscription with id: %s is notifying but not on db", subscriptionID), tags)
		return nil, nil
	} else if err != nil {
		//Sentry already got this error
		return nil, err
	}
	calendar, err = s.database.RetrieveCalendarFromSubscription(subscriptionID)
	if err != nil {
		s.sentry.CaptureErrorAndWait(err, tags)
		log.Errorf("error refreshing outlook account")
		return nil, err
	}
	recoveredPanic, sentryID := s.sentry.CapturePanicAndWait(func() {
		err = calendar.GetAccount().Refresh()
	}, tags)

	if recoveredPanic != nil {
		log.Errorf("panic recovered with sentry ID: %s", sentryID)
		return nil, fmt.Errorf("panic was launched")
	}
	if err != nil {
		s.sentry.CaptureErrorAndWait(err, tags)
		log.Errorf("error refreshing outlook account")
		return nil, err
	}
	return calendar, nil
}
