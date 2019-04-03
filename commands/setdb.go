package commands

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func createTables(db *sql.DB) {
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "active" (
		"date" TEXT,
		"source" TEXT,
		"num" TEXT,
		"cat" TEXT,
		"desc" TEXT,
		"link" TEXT,
		"contact" TEXT,
		"price" TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()
	stmt.Close()
	log.Info("created table active")

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS "archive" (
		"date" TEXT,
		"source" TEXT,
		"num" TEXT UNIQUE ON CONFLICT IGNORE,
		"cat" TEXT,
		"desc" TEXT,
		"link" TEXT,
		"contact" TEXT,
		"price" TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()
	stmt.Close()
	log.Info("created table archive")
}

// SetDefaultDB creates database.
// if test == true test.db is created
// else adverts.db is created
func SetDefaultDB(test bool) {
	if test {
		db, err := sql.Open("sqlite3", "./db/test.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		createTables(db)
		log.Info("setting up TEST database finished")
	} else {
		db, err := sql.Open("sqlite3", "./db/adverts.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		createTables(db)
		log.Info("setting up database finished")
	}
}
