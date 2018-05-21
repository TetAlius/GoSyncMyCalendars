package db

import (
	"database/sql"
	"net/http"

	"github.com/TetAlius/GoSyncMyCalendars/api"
	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
	"github.com/google/uuid"
)

func (data Database) RetrieveAllSubscriptions(principalSubscriptionUUID string) (subscriptions []api.SubscriptionManager, err error) {
	subscription, err := data.getSubscription(principalSubscriptionUUID)
	if err != nil {
		log.Errorf("some error retrieving principal subscription: %s", err.Error())
		return
	}
	subscriptions = append(subscriptions, subscription)
	rows, err := data.DB.Query("select s2.uuid from subscriptions as s2 join calendars c2 on s2.calendar_uuid = c2.uuid where c2.parent_calendar_uuid IN (select s.calendar_uuid from subscriptions as s where s.uuid=$1)", principalSubscriptionUUID)
	defer rows.Close()
	for rows.Next() {
		var uid uuid.UUID
		rows.Scan(&uid)
		subscription, err := data.getSubscription(uid.String())
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
func (data Database) getSubscription(subscriptionUUID string) (subscription api.SubscriptionManager, err error) {
	var id string
	var uid uuid.UUID
	var calendarUUID string
	var kind int
	err = data.DB.QueryRow("select subscriptions.id,subscriptions.uuid,subscriptions.calendar_uuid, a.kind from subscriptions join calendars c2 on subscriptions.calendar_uuid = c2.uuid join accounts a on c2.account_email = a.email where subscriptions.uuid = $1", subscriptionUUID).Scan(&id, &uid, &calendarUUID, &kind)
	switch {
	case err == sql.ErrNoRows:
		log.Debugf("no subscription with that uuid: %s.", subscriptionUUID)
		return nil, &customErrors.NotFoundError{Code: http.StatusNotFound}
	case err != nil:
		log.Errorf("error looking for subscription with uuid: %s", subscriptionUUID)
		return
	}
	switch kind {
	case api.GOOGLE:
		subscription = &api.GoogleSubscription{ID: id, Uuid: uid}
	case api.OUTLOOK:
		subscription = &api.OutlookSubscription{ID: id, Uuid: uid}
	default:
		return nil, &WrongKindError{subscriptionUUID}
	}
	return
}
