package calendar

import (
	"database/sql"

	"github.com/TetAlius/GoSyncMyCalendars/customErrors"
	"github.com/TetAlius/GoSyncMyCalendars/database"
	log "github.com/TetAlius/GoSyncMyCalendars/logger"
)

// Create creates a new calendar with the account given, the id of the calendar
// on the server and the mail associated with the calendar.
// mail must be registered before on the accounts table.
// If the calendar could not be created it will raise a RowsError with the information given
func Create(db *sql.DB, account string, id string, mail string) (err error) {
	//Prepare database unless database is already given
	if db == nil {
		db = database.OpenDB()
		//Close database once all is finished, unless database is already given
		//because it will be on a for loop and will be closed at the end
		defer db.Close()
	}

	sentence, err := getSentenceToCalendarInsert(account)
	if err != nil {
		log.Errorf("GetSentenceToCalendarInsert unexpectedly failed: %v", err)
	}

	result, err := db.Exec(sentence, id)
	if err != nil {
		log.Errorf("db.Exec unexpectedly failed: %v", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Errorf("result.RowsAffected unexpectedly failed: %v", err)
	}
	if n != 1 {
		return customErrors.RowsError{Expected: 1, Received: n, Message: "On calendar.Create with id: " + id + " mail: " + mail + " account: " + account}
	}
	return

}

//Update TODO
func Update(db *sql.DB, account string, id string, info interface{}) (err error) {
	//Prepare database unless database is already given
	if db == nil {
		db = database.OpenDB()
		//Close database once all is finished, unless database is already given
		//because it will be on a for loop and will be closed at the end
		defer db.Close()
	}
	return

}

//Get TODO
func Get(db *sql.DB, account string, id string) (v interface{}, err error) {
	//Prepare database unless database is already given
	if db == nil {
		db = database.OpenDB()
		//Close database once all is finished, unless database is already given
		//because it will be on a for loop and will be closed at the end
		defer db.Close()
	}

	sentence, err := getSentenceToCalendarRead(account)
	if err != nil {
		log.Errorf("GetSentenceToCalendarRead unexpectedly failed: %v", err)
		return
	}

	rows, err := db.Query(sentence, id)
	if err != nil {
		log.Errorf("db.Exec unexpectedly failed: %v", err)
		return
	}

	v, err = database.HandleCalendarRowsForAccount(account, rows)

	return

}

// Delete deletes a calendar that is registered on the database
// if the id given does not correspond to a calendar, it will raise a RowsError
// whith the id of the calendar
func Delete(db *sql.DB, id string) (err error) {
	//Prepare database unless database is already given
	if db == nil {
		db = database.OpenDB()
		//Close database once all is finished, unless database is already given
		//because it will be on a for loop and will be closed at the end
		defer db.Close()
	}

	sentence := "DELETE FROM calendars WHERE id = $1"

	result, err := db.Exec(sentence, id)
	if err != nil {
		log.Errorf("db.Exec unexpectedly failed: %v", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		log.Errorf("result.RowsAffected unexpectedly failed: %v", err)
	}
	if n != 1 {
		return customErrors.RowsError{Expected: 1, Received: n, Message: "On calendar.Delete with id: " + id}
	}
	return
}
