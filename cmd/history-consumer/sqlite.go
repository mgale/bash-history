package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mgale/bash-history.git/internal/events"
	_ "modernc.org/sqlite"
)

func createDBFileIfNotExists(db_filename string) error {
	if _, err := os.Stat(db_filename); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Println("Creating DB file: ", db_filename)
		f, err := os.Create(db_filename)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	return nil
}

func createDBTablesIfNotExists(db *sql.DB, myTable string) error {
	query, err := db.Prepare(myTable)
	if err != nil {
		return err
	}
	result, err := query.Exec()
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("Tables created:", rows)
	return nil
}

func createDBTableIndexesIfNotExists(db *sql.DB, indexes []string) error {
	for _, index := range indexes {
		query, err := db.Prepare(index)
		if err != nil {
			return err
		}
		result, err := query.Exec()
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		fmt.Println("Indexes created:", rows)
	}
	return nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS bashhistory (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			timestamp integer(4) not null default (strftime('%s','now')),
			username text not null,
			command text not null,
			UNIQUE(username, command));`,
		`CREATE TABLE IF NOT EXISTS bashsessions (
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				timestamp integer(4) not null default (strftime('%s','now')),
				pid integer(8) not null,
				username text not null,
				command text not null);`,
	}

	for _, table := range tables {
		err := createDBTablesIfNotExists(db, table)
		if err != nil {
			return err
		}
	}

	return nil
}

func createTableIndexes(db *sql.DB) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_bashhistory_timestamp ON bashhistory(timestamp);",
	}

	return createDBTableIndexesIfNotExists(db, indexes)
}

/* Note that any field not in the insert list will be set to NULL if the row already exists in the table.
This is why there's a subselect for the id column: In the replacement case the statement would set it
to NULL and then a fresh ID would be allocated.
This approach can also be used if you want to leave particular field values alone if the row in the
replacement case but set the field to NULL in the insert case.
*/
func insertHistoryTableEvent(event events.ReadEvent, db *sql.DB) (int64, error) {

	insert, err := db.Prepare("REPLACE INTO bashhistory(id, username, command) VALUES((SELECT id from bashhistory WHERE username=? AND command = ?),?,?)")
	if err != nil {
		return 0, err
	}
	result, err := insert.Exec(event.Username, event.Line, event.Username, event.Line)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func insertSessionTableEvent(event events.ReadEvent, db *sql.DB) error {
	insert, err := db.Prepare("INSERT INTO bashsessions(pid, username, command) VALUES(?,?,?)")
	if err != nil {
		return err
	}
	_, err = insert.Exec(event.Pid, event.Username, event.Line)
	if err != nil {
		return err
	}
	return nil
}

func createDBConnection(db_filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", db_filename)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func handleEvents(ctx context.Context, streamEventsChannel chan events.ReadEvent, docEventsChannel chan events.DocumentEvent, db *sql.DB, verbose bool) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-streamEventsChannel:
			if verbose {
				log.Printf("Received event: %+v", event)
			}
			id, err := insertHistoryTableEvent(event, db)
			if err != nil {
				return err
			}
			err = insertSessionTableEvent(event, db)
			if err != nil {
				return err
			}
			docEvent := events.DocumentEvent{
				ID:        id,
				Username:  event.Username,
				Command:   event.Line,
				Timestamp: time.Now().Unix(),
			}
			log.Printf("Sending event: %+v", docEvent)
			docEventsChannel <- docEvent
		}
	}
}
