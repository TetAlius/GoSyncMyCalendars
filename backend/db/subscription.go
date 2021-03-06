package db

import (
	"database/sql"
	"errors"

	"fmt"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
)

// Returns all subscriptions from a user
func (data Database) RetrieveAllSubscriptionsFromUser(principalSubscriptionUUID string, userEmail string, userUUID string) (subscriptions []api.SubscriptionManager, err error) {
	subscription, err := data.getSubscription(principalSubscriptionUUID, userEmail, userUUID)
	if err != nil {
		log.Errorf("some error retrieving principal subscription: %s", err.Error())
		return
	}
	subscriptions = append(subscriptions, subscription)
	rows, err := data.client.Query("select s2.uuid from subscriptions as s2 join calendars c2 on s2.calendar_uuid = c2.uuid join accounts a on c2.account_email = a.email join users u on a.user_uuid = u.uuid where u.uuid = $1 and u.email = $2 and c2.parent_calendar_uuid IN (select s.calendar_uuid from subscriptions as s where s.uuid=$3)", userUUID, userEmail, principalSubscriptionUUID)
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error querying principal subscription with uuid: %s", principalSubscriptionUUID)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var uid uuid.UUID
		err := rows.Scan(&uid)
		if err != nil {
			log.Errorf("error scanning query: %s", err.Error())
			break
		}
		subscription, err := data.getSubscription(uid.String(), userEmail, userUUID)
		if err != nil {
			log.Errorf("error getting subscription with uuid: %s", uid)
			break
		} else {
			subscriptions = append(subscriptions, subscription)
		}
	}
	if err != nil {
		log.Errorf("some error retrieving all subscriptions: %s", err.Error())
		return
	}

	return
}

// Deletes a subscription
func (data Database) deleteSubscription(transaction *sql.Tx, subscription api.SubscriptionManager) (err error) {
	stmt, err := transaction.Prepare("delete from subscriptions where subscriptions.uuid = $1")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing sql: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(subscription.GetUUID())
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error executing delete: %s", err.Error())
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error getting affected rows: %s", err.Error())
		return
	}
	if affect != 1 {
		err = fmt.Errorf("more than one row affected by deletion uuid: %s", subscription.GetUUID())
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		transaction.Rollback()
		return
	}
	return
}

// Retrieves all subscription that are expiring or are expired
func (data Database) GetExpiredSubscriptions() (subscriptions []api.SubscriptionManager, err error) {
	rows, err := data.client.Query("select subscriptions.uuid, u.email, u.uuid from subscriptions join calendars c2 on subscriptions.calendar_uuid = c2.uuid join accounts a on c2.account_email = a.email join users u on a.user_uuid = u.uuid where subscriptions.expiration_date <= current_date")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error retrieving all subscription expired: %s", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var subscription api.SubscriptionManager
		var subscriptionUUID string
		var userEmail string
		var userUUID string
		err = rows.Scan(&subscriptionUUID, &userEmail, &userUUID)
		if err != nil {
			data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
			log.Errorf("error scanning results: %s", err.Error())
			return nil, err
		}
		subscription, err = data.getSubscription(subscriptionUUID, userEmail, userUUID)
		if err != nil {
			log.Errorf("error getting subscription with uuid: %s email: %s user uuid: %s error: %s", subscriptionUUID, userEmail, userUUID, err.Error())
			err = nil
			continue
		}
		subscriptions = append(subscriptions, subscription)
	}
	return

}

// Method that updates the info of a subscription
func (data Database) UpdateSubscription(subscription api.SubscriptionManager) (err error) {
	stmt, err := data.client.Prepare("update subscriptions set id = $1, type = $2, expiration_date = $3, resource_id = $4 where uuid = $5")
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error preparing query: %s", err.Error())
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(subscription.GetID(), subscription.GetType(), subscription.GetExpirationDate(), subscription.GetResourceID(), subscription.GetUUID())
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error executing query: %s", err.Error())
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error retrieving rows affected: %s", err.Error())
		return
	}
	if affect != 1 {
		err = errors.New(fmt.Sprintf("could not update subscription with uuid: %s", subscription.GetUUID()))
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		return
	}
	return
}

// Method that retrieves a subscription
func (data Database) getSubscription(subscriptionUUID string, userEmail string, userUUID string) (subscription api.SubscriptionManager, err error) {
	var id string
	var uid uuid.UUID
	var calendarUUID string
	var kind int
	var typ string
	var resourceID string
	err = data.client.QueryRow("select subscriptions.id,subscriptions.uuid,subscriptions.calendar_uuid, a.kind, subscriptions.type, subscriptions.resource_id from subscriptions join calendars c2 on subscriptions.calendar_uuid = c2.uuid join accounts a on c2.account_email = a.email join users u on a.user_uuid = u.uuid where subscriptions.uuid = $1 and u.uuid = $2 and u.email = $3", subscriptionUUID, userUUID, userEmail).Scan(&id, &uid, &calendarUUID, &kind, &typ, &resourceID)
	switch {
	case err == sql.ErrNoRows:
		err = &customErrors.NotFoundError{Message: fmt.Sprintf("no subscription with that uuid: %s.", subscriptionUUID)}
		log.Debugf("no subscription with that uuid: %s.", subscriptionUUID)
		return nil, err
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Errorf("error looking for subscription with uuid: %s", subscriptionUUID)
		return
	}
	calendar, _ := data.findCalendarFromUser(userEmail, userUUID, calendarUUID)
	switch kind {
	case api.GOOGLE:
		subscription = api.RetrieveGoogleSubscription(id, uid, calendar, resourceID)
	case api.OUTLOOK:
		subscription = api.RetrieveOutlookSubscription(id, uid, calendar, typ)
	default:
		data.sentry.CaptureErrorAndWait(&customErrors.WrongKindError{Mail: subscriptionUUID}, map[string]string{"database": "backend"})
		return nil, &customErrors.WrongKindError{Mail: subscriptionUUID}
	}
	return
}

// Method that returns if a subscription is stored on DB
func (data Database) ExistsSubscriptionFromID(ID string) (ok bool, err error) {
	err = data.client.QueryRow("SELECT true FROM subscriptions where subscriptions.id = $1", ID).Scan(&ok)
	switch {
	case err == sql.ErrNoRows:
		err = &customErrors.NotFoundError{Message: fmt.Sprintf("No subscription with id: %s", ID)}
		log.Debugf("No subscription with id: %s", ID)
		return true, err
	case err != nil:
		data.sentry.CaptureErrorAndWait(err, map[string]string{"database": "backend"})
		log.Debugf("error looking for subscription with id: %s", ID)
		return false, err
	}
	return
}
