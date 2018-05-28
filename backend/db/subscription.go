package db

import (
	"database/sql"
	"net/http"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
)

func (data Database) RetrieveAllSubscriptionsFromUser(principalSubscriptionUUID string, userEmail string, userUUID string) (subscriptions []api.SubscriptionManager, err error) {
	subscription, err := data.getSubscription(principalSubscriptionUUID, userEmail, userUUID)
	if err != nil {
		log.Errorf("some error retrieving principal subscription: %s", err.Error())
		return
	}
	subscriptions = append(subscriptions, subscription)
	rows, err := data.DB.Query("select s2.uuid from subscriptions as s2 join calendars c2 on s2.calendar_uuid = c2.uuid join accounts a on c2.account_email = a.email join users u on a.user_uuid = u.uuid where u.uuid = $1 and u.email = $2 and c2.parent_calendar_uuid IN (select s.calendar_uuid from subscriptions as s where s.uuid=$3)", userUUID, userEmail, principalSubscriptionUUID)
	defer rows.Close()
	for rows.Next() {
		var uid uuid.UUID
		rows.Scan(&uid)
		subscription, err := data.getSubscription(uid.String(), userEmail, userUUID)
		if err != nil {
			log.Errorf("error getting subscription with uuid: %s", uid)
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

func (data Database) DeleteSubscription(subscription api.SubscriptionManager) (err error) {
	stmt, err := data.DB.Prepare("delete from subscriptions where subscriptions.uuid = $1")
	defer stmt.Close()
	if err != nil {
		log.Errorf("error preparing sql: %s", err.Error())
		return
	}
	res, err := stmt.Exec(subscription.GetUUID())
	if err != nil {
		log.Errorf("error executing delete: %s", err.Error())
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Errorf("error getting affected rows: %s", err.Error())
		return
	}
	log.Debugf("affected %d rows", affect)
	return
}

func (data Database) GetExpiredSubscriptions() (subscriptions []api.SubscriptionManager, err error) {
	rows, err := data.DB.Query("select subscriptions.uuid, u.email, u.uuid from subscriptions join calendars c2 on subscriptions.calendar_uuid = c2.uuid join accounts a on c2.account_email = a.email join users u on a.user_uuid = u.uuid where subscriptions.expiration_date <= current_date")
	if err != nil {
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
			log.Errorf("error scanning results: %s", err.Error())
			return nil, err
		}
		subscription, err = data.getSubscription(subscriptionUUID, userEmail, userUUID)
		if err != nil {
			log.Errorf("error getting subscription with uuid: %s email: %s user uuid: %s error: %s", subscriptionUUID, userEmail, userUUID, err.Error())
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return

}

func (data Database) UpdateSubscription(subscription api.SubscriptionManager) (err error) {

	return
}

func (data Database) getSubscription(subscriptionUUID string, userEmail string, userUUID string) (subscription api.SubscriptionManager, err error) {
	var id string
	var uid uuid.UUID
	var calendarUUID string
	var kind int
	var typ string
	err = data.DB.QueryRow("select subscriptions.id,subscriptions.uuid,subscriptions.calendar_uuid, a.kind, subscriptions.type from subscriptions join calendars c2 on subscriptions.calendar_uuid = c2.uuid join accounts a on c2.account_email = a.email join users u on a.user_uuid = u.uuid where subscriptions.uuid = $1 and u.uuid = $2 and u.email = $3", subscriptionUUID, userUUID, userEmail).Scan(&id, &uid, &calendarUUID, &kind, &typ)
	switch {
	case err == sql.ErrNoRows:
		log.Debugf("no subscription with that uuid: %s.", subscriptionUUID)
		return nil, &customErrors.NotFoundError{Code: http.StatusNotFound}
	case err != nil:
		log.Errorf("error looking for subscription with uuid: %s", subscriptionUUID)
		return
	}
	calendar, _ := data.findCalendarFromUser(userEmail, userUUID, calendarUUID)
	switch kind {
	case api.GOOGLE:
		subscription = api.RetrieveGoogleSubscription(id, uid, calendar)
	case api.OUTLOOK:
		subscription = api.RetrieveOutlookSubscription(id, uid, calendar, typ)
	default:
		return nil, &customErrors.WrongKindError{Mail: subscriptionUUID}
	}
	return
}
