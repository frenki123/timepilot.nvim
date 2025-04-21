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
	schema := `
    CREATE TABLE IF NOT EXISTS project_timer (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        path TEXT NOT NULL,
        started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        end_at DATETIME
    );
    CREATE TABLE IF NOT EXISTS file_timer (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        project_id INTEGER NOT NULL,
        path TEXT NOT NULL,
        filetype TEXT,
        started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        end_at DATETIME,
        FOREIGN KEY(project_id) REFERENCES project_timer(id)
    );
    `
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}
