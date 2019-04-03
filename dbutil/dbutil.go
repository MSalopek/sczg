package dbutil

import (
	"database/sql"
	"fmt"

	"sczg/web"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type transactionError struct {
	err string
	fn  string
}

// NewErr returns new transactionError
func NewErr(e string, fn string) error {
	return &transactionError{e, fn}
}

func errWithLog(e string, fn string) error {
	tErr := NewErr(e, fn)
	log.Errorf(tErr.Error())
	return tErr
}

func (t *transactionError) Error() string {
	return fmt.Sprintf("error executing transaction in {fn: %v}  details {%v}", t.fn, t.err)
}

// Storage is a database wrapper
type Storage struct {
	*sql.DB
}

// InitStorage established DB Conn
func InitStorage(db string) (*Storage, error) {
	base, err := sql.Open("sqlite3", db)
	if err != nil {
		log.Printf("failed to initialize sqlite {%v}\n", err)
		return nil, err
	}
	return &Storage{base}, nil
}

// insert handles filtering and caching of new ads into "active" table
func (s *Storage) insert(ads []web.Advert, filter map[string]struct{}) error {
	tx, err := s.Begin()
	if err != nil {
		return errWithLog(err.Error(), "InsAds")
	}
	stmt, err := tx.Prepare("INSERT INTO active VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errWithLog(err.Error(), "InsAds")
	}
	defer stmt.Close()
	insCtr := 0
	for _, v := range ads {
		if _, ok := filter[v.AdNum]; !ok { // HACK: avoids using sort and search, uses map lookup instead
			_, err := stmt.Exec(v.Date, v.Source, v.AdNum, v.Cat, v.Desc, v.Link, v.Contact, v.Rate)
			if err != nil {
				tx.Rollback()
				return errWithLog(err.Error(), "insert")
			}
			insCtr++
		}
	}
	err = tx.Commit()
	if err != nil {
		return errWithLog(err.Error(), "insert")
	}
	log.Infof("inserted {%v} out of {%v} found entries", insCtr, len(ads))
	return nil
}

// getCurEntryNums fetches ad numbers of currently
// active ads and returns the result as map.
func (s *Storage) getCurEntryNums() (map[string]struct{}, error) {
	var entryNum string
	curEntries := make(map[string]struct{})
	rows, err := s.Query("SELECT num from active")
	if err != nil {
		return map[string]struct{}{}, errWithLog(err.Error(), "getCurEntryNums")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&entryNum)
		if err != nil {
			if err == sql.ErrNoRows {
				// table is empty just return empty map
				return map[string]struct{}{}, nil
			}
			return nil, errWithLog(err.Error(), "getCurEntryNums")
		}
		curEntries[entryNum] = struct{}{}
	}
	return curEntries, nil
}

// InsertNewAds filters adverts and inserts new to DB.
func (s *Storage) InsertNewAds(ads []web.Advert) error {
	filter, err := s.getCurEntryNums()
	if err != nil {
		return errWithLog(err.Error(), "InsertNewAds")
	}
	if err := s.insert(ads, filter); err != nil {
		return errWithLog(err.Error(), "InsertNewAds")
	}
	return nil
}

// PurgeOutdatedActive deletes entries
// under ACTIVE older than datetime(now - param::days)
func (s *Storage) PurgeOutdatedActive(days int) error {
	tx, _ := s.Begin()
	prep := fmt.Sprintf("DELETE FROM active WHERE date < datetime('now', '-%v days','localtime')", days)
	stmt, err := tx.Prepare(prep)
	if err != nil {
		return errWithLog(err.Error(), "PurgeOutdatedActive")
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return errWithLog(err.Error(), "PurgeOutdatedActive")
	}
	err = tx.Commit()
	if err != nil {
		return errWithLog(err.Error(), "PurgeOutdatedActive")
	}
	log.Info("purging old entries completed")
	return nil
}

// FetchAllActive querys db and returns all ACTIVE adverts
// as k:v pairs where k == category name && v == []web.Advert
func (s *Storage) FetchAllActive() ([]web.Advert, error) {
	var ad web.Advert
	var allEntries []web.Advert
	rows, err := s.Query("SELECT * FROM active")
	if err != nil {
		return nil, errWithLog(err.Error(), "FetchAllActive")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&ad.Date, &ad.Source, &ad.AdNum, &ad.Cat, &ad.Desc, &ad.Link, &ad.Contact, &ad.Rate)
		if err != nil {
			if err == sql.ErrNoRows {
				// table is empty just return empty map
				return []web.Advert{}, nil
			}
			return nil, errWithLog(err.Error(), "FetchAllActive")
		}
		allEntries = append(allEntries, ad)
	}
	return allEntries, nil
}

// FetchActiveByCategory queries db for ACTIVE ads of a specific category
func (s *Storage) FetchActiveByCategory(cat string) ([]web.Advert, error) {
	var ad web.Advert
	var entries []web.Advert
	qstr := fmt.Sprintf("SELECT * FROM active WHERE cat = '%s'", cat)
	rows, err := s.Query(qstr)
	if err != nil {
		return nil, errWithLog(err.Error(), "FetchActiveByCategory")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&ad.Date, &ad.Source, &ad.AdNum, &ad.Cat, &ad.Desc, &ad.Link, &ad.Contact, &ad.Rate)
		if err != nil {
			if err == sql.ErrNoRows {
				// table is empty just return empty map
				return []web.Advert{}, nil
			}
			return nil, errWithLog(err.Error(), "FetchActiveByCategory")
		}
		entries = append(entries, ad)
	}
	return entries, nil
}

// FetchFreshActive fetches ads that are < 1 day old
func (s *Storage) FetchFreshActive() ([]web.Advert, error) {
	var ad web.Advert
	var freshAds []web.Advert
	rows, err := s.Query("SELECT * FROM active WHERE date >= datetime('now', 'start of day')")
	if err != nil {
		return nil, errWithLog(err.Error(), "FetchFreshActive")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&ad.Date, &ad.Source, &ad.AdNum, &ad.Cat, &ad.Desc, &ad.Link, &ad.Contact, &ad.Rate)
		if err != nil {
			if err == sql.ErrNoRows {
				// table is empty just return empty map
				return []web.Advert{}, nil
			}
			return nil, errWithLog(err.Error(), "FetchFreshActive")
		}
		freshAds = append(freshAds, ad)
	}
	return freshAds, nil
}

// ArchiveEntries inserts entries from "active" into "archive" table.
// Param days defines minimum age of ads to be archived.
func (s *Storage) ArchiveEntries(days int) error {
	tx, _ := s.Begin()
	prep := fmt.Sprintf("INSERT OR IGNORE INTO archive SELECT * FROM active WHERE date > datetime('now','-%v days', 'start of day')", days)
	stmt, err := tx.Prepare(prep)
	if err != nil {
		tx.Rollback()
		return errWithLog(err.Error(), "ArchiveEntries")
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return errWithLog(err.Error(), "ArchiveEntries")
	}
	if err = tx.Commit(); err != nil {
		return errWithLog(err.Error(), "ArchiveEntries")
	}
	log.Infof("finished archiving ads older than %v days", days)
	return nil
}
