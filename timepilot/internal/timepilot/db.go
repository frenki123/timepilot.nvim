package timepilot

import (
	"os"
	"path"

	"github.com/jmoiron/sqlx"
)

func GetDB(directory string) (*sqlx.DB, error) {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, err
	}
	uri := path.Join(directory, "timepilot.db")
	db, err := sqlx.Connect("sqlite", uri)
	if err != nil {
		return nil, err
	}
	if err := CreateTable(db); err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sqlx.DB) error {
	schema := `CREATE TABLE IF NOT EXISTS actions (
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        filename TEXT,
        filetype TEXT,
        action TEXT,
        date DATE DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}
