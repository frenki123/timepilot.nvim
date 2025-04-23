package timepilot

import (
	"encoding/json"
	"errors"
	"strings"
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

type Response struct {
	Type string `json:"type"`
	Data any    `json:"data"`
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
	return Response{
		Type: "DEBUG",
		Data: "Stored Project Enter",
	}, nil
}

func storeSessionEnter(db *sqlx.DB, project Project) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmnt := `INSERT INTO project_timer (path) VALUES (?) RETURNING id`
	res, err := tx.Exec(stmnt, strings.Trim(project.Filepath, "\n"))
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
	return Response{
		Type: "DEBUG",
		Data: "Stored Project Leave",
	}, nil
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
	_, err = tx.Exec(stmnt, activeProjectId, strings.Trim(project.Filepath, "\n"))
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
	if buffer.Filepath == "" || buffer.Filetype == "" || buffer.Project == "" {
		return "", errors.New("Non valid params")
	}

	if err := storeBufferEnter(db, buffer); err != nil {
		return "", err
	}
	return Response{
		Type: "DEBUG",
		Data: "Stored Buffer Enter",
	}, nil
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
	if buffer.Filepath == "" || buffer.Project == "" {
		return "", errors.New("Non valid params")
	}

	if err := updateBufferLeave(db, buffer); err != nil {
		return "", err
	}
	return Response{
		Type: "DEBUG",
		Data: "Stored Buffer Leave",
	}, nil
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

type ProjectName struct {
	Filepath string `json:"project_name"`
}

func GetProjectTime(db *sqlx.DB, params json.RawMessage) (any, error) {
	var project ProjectName
	err := json.Unmarshal(params, &project)
	if err != nil {
		return "", err
	}
	if project.Filepath == "" {
		return "", errors.New("Non valid params")
	}
	time, err := getProjectTime(db, project)
	if err != nil {
		return "", err
	}
	return Response{
		Type: "INFO/TIME",
		Data: time,
	}, nil
}

func getProjectTime(db *sqlx.DB, proj ProjectName) (int, error) {
	stmnt := `
    SELECT
        SUM(strftime('%s', end_at) - strftime('%s', started_at))/60 AS time
    FROM
        project_timer
    WHERE
        end_at IS NOT NULL
        AND path=?
    GROUP BY
        path;
    `
	row := db.QueryRow(stmnt, strings.Trim(proj.Filepath, "\n"))
	var time int
	err := row.Scan(&time)
	if err != nil {
		return 0, err
	}
	return time, nil
}

type FileResult struct {
	Filepath string `json:"filepath"`
	Filetype string `json:"filetype"`
	Time     int    `json:"time"`
}

func GetMostEditedFile(db *sqlx.DB, params json.RawMessage) (any, error) {
	var project ProjectName
	err := json.Unmarshal(params, &project)
	if err != nil {
		return "", err
	}
	if project.Filepath == "" {
		return "", errors.New("Non valid params")
	}
	file, err := getMostEditedFile(db, project)
	if err != nil {
		return "", err
	}
	return Response{
		Type: "INFO/FILE",
		Data: file,
	}, nil
}

func getMostEditedFile(db *sqlx.DB, proj ProjectName) (FileResult, error) {
	stmnt := `
    SELECT
        f.path AS filepath,
        f.filetype as filetype,
        SUM(strftime('%s', f.end_at) - strftime('%s', f.started_at))/60 AS time
    FROM
        file_timer f
    JOIN
        project_timer p ON f.project_id = p.id
    WHERE
        f.end_at IS NOT NULL
        AND p.path=?
    GROUP BY
        f.path
    ORDER BY time DESC;
    `
	var res FileResult
	err := db.Get(&res, stmnt, strings.Trim(proj.Filepath, "\n"))
	if err != nil {
		return FileResult{}, err
	}
	return res, nil
}
