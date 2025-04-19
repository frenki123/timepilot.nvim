package timepilot

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type ActionParams struct {
	Filename string
	Filetype string
	Action   string
}

type SaveStore struct {
	Filename string
	Filetype string
	Action   string
	Time     time.Time
}

func ActionControler(db *sqlx.DB, params json.RawMessage) (any, error) {
    var action ActionParams
    err := json.Unmarshal(params, &action)
    if  err != nil {
        return "", err
    }
    if action.Action == "" || action.Filename == "" || action.Filetype == "" {
        return "", errors.New("Non valid params")
    }
    if err := StoreActionToDB(db, action); err != nil {
        return "", err
    }
	return "Stored", nil
}

func StoreActionToDB(db *sqlx.DB, action ActionParams) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    stmnt := `INSERT INTO actions (filename, filetype, action) VALUES (?,?,?)`
    _, err = tx.Exec(stmnt, action.Filename, action.Filetype, action.Action)
    if err != nil {
        tx.Rollback()
        return err
    }
    if err := tx.Commit(); err != nil {
        return err
    }
    return nil
}
