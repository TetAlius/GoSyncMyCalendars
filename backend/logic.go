package backend

import (
	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

func (s *Server) manageSynchronizationOutlook(notifications []api.OutlookSubscriptionNotification) (err error) {
	for _, subscription := range notifications {
		err = s.manageSubscription(subscription.SubscriptionID, subscription.Data.ID, map[string]string{"sync": "outlook"})
		if err != nil {
			log.Errorf("error managing outlook subscription ID: %s", subscription.SubscriptionID)
			return err
		}
	}
	return err
}

func (s *Server) manageSynchronizationGoogle(subscriptionID string, eventID string) (err error) {
	err = s.manageSubscription(subscriptionID, eventID, map[string]string{"sync": "google"})
	if err != nil {
		log.Errorf("error managing google subscription ID: %s", subscriptionID)
		return err
	}
	return err
}

func (s *Server) manageSubscription(subscriptionID string, eventID string, tags map[string]string) (err error) {
	ok, err := s.database.ExistsSubscriptionFromID(subscriptionID)
	if err != nil && ok {
		//Ignore this subscription
		//Perhaps something went wrong when deleting subscription...
		//s.sentry.CaptureMessageAndWait(fmt.Sprintf("outlook subscription with id: %s is notifying but not on db", subscriptionID), tags)
		return nil
	} else if err != nil {
		//Sentry already got this error
		return err
	}
	calendar, err := s.database.RetrieveCalendarFromSubscription(subscriptionID)

	if err != nil {
		s.sentry.CaptureErrorAndWait(err, tags)
		log.Errorf("error refreshing outlook account")
		return err
	}
	recoveredPanic, sentryID := s.sentry.CapturePanicAndWait(func() {
		err = calendar.GetAccount().Refresh()
	}, tags)

	if recoveredPanic != nil {
		log.Errorf("panic recovered with sentry ID: %s", sentryID)
		return fmt.Errorf("panic was launched")
	}
	if err != nil {
		s.sentry.CaptureErrorAndWait(err, tags)
		log.Errorf("error refreshing outlook account")
		return err
	}
	//TODO: update account
	//go func() { s.database.UpdateAccount(calendar.GetAccount()) }()
	onCloud := true
	event, err := calendar.GetEvent(eventID)
	if _, ok := err.(*customErrors.NotFoundError); ok {
		onCloud = false
		err = nil
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
	if onCloud && onDB && s.database.EventAlreadyUpdated(event) {
		return nil
	}
	//TODO: manage this error, if returns none event maybe because is deleted

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
