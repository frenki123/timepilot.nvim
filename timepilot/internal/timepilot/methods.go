package timepilot

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var activeProjectId int64

type Project struct {
	Filepath string `json:"project"`
}
type Session struct {
	Id string `json:"sessionId"`
}
type Buffer struct {
	Project  string
	Filepath string `json:"filename"`
	Filetype string
}

type SaveStore struct {
	Filename string
	Filetype string
	Action   string
	Time     time.Time
}

func SessionEnter(db *sqlx.DB, params json.RawMessage) (any, error) {
	var project Project
	err := json.Unmarshal(params, &project)
	if err != nil {
		return "", err
	}
	if project.Filepath == "" {
		return "", errors.New("Non valid params")
	}
	if err := storeSessionEnter(db, project); err != nil {
		return "", err
	}
	return "Stored", nil
}

func storeSessionEnter(db *sqlx.DB, project Project) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmnt := `INSERT INTO project_timer (path) VALUES (?) RETURNING id`
	res, err := tx.Exec(stmnt, project.Filepath)
	activeProjectId, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func SessionLeave(db *sqlx.DB, params json.RawMessage) (any, error) {
	var project Project
	err := json.Unmarshal(params, &project)
	if err != nil {
		return "", err
	}
	if project.Filepath == "" {
		return "", errors.New("Non valid params")
	}
	if err := updatedSessionLeave(db, project); err != nil {
		return "", err
	}
	return "Stored", nil
}

func updatedSessionLeave(db *sqlx.DB, project Project) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmnt := `
    UPDATE project_timer
    SET end_at=CURRENT_TIMESTAMP
    WHERE ID = ? AND path = ?;
    `
	_, err = tx.Exec(stmnt, activeProjectId, project.Filepath)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func BufferEnter(db *sqlx.DB, params json.RawMessage) (any, error) {
	var buffer Buffer
	err := json.Unmarshal(params, &buffer)
	if err != nil {
		return "", err
	}
	//if buffer.Filepath == "" || buffer.Filetype == "" || buffer.Project == "" {
	//return "", errors.New("Non valid params")
	//}

	if err := storeBufferEnter(db, buffer); err != nil {
		return "", err
	}
	return "Stored", nil
}

func storeBufferEnter(db *sqlx.DB, buffer Buffer) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmnt := `
    INSERT INTO file_timer 
    (project_id, path, filetype) VALUES (?, ?, ?);
    `
	_, err = tx.Exec(stmnt, activeProjectId, buffer.Filepath, buffer.Filetype)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func BufferLeave(db *sqlx.DB, params json.RawMessage) (any, error) {
	var buffer Buffer
	err := json.Unmarshal(params, &buffer)
	if err != nil {
		return "", err
	}
	// if buffer.Filepath == "" || buffer.Project == "" {
	// 	return "", errors.New("Non valid params")
	// }

	if err := updateBufferLeave(db, buffer); err != nil {
		return "", err
	}
	return "Stored", nil
}

func updateBufferLeave(db *sqlx.DB, buffer Buffer) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmnt := `
    UPDATE file_timer
    SET end_at=CURRENT_TIMESTAMP
    WHERE path = ? AND end_at IS NULL;
    `
	_, err = tx.Exec(stmnt, buffer.Filepath)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
