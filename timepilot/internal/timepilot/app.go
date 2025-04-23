package timepilot

import (
	"encoding/json"
	"log"
	"os"

	"github.com/frenki123/timepilot.nvim/timepilot/internal/tprpc"
	"github.com/jmoiron/sqlx"
)

type App struct {
	srv    tprpc.Server
	db     *sqlx.DB
	config Config
}

type Config struct {
	DBPath string
}

func InitConfig() Config {
	path := os.Getenv("HOME") + "/.config/timepilot"
	return Config{
		DBPath: path,
	}
}

func NewApp() App {
	config := InitConfig()
	db, err := GetDB(config.DBPath)
	if err != nil {
		log.Fatalf("DB Error: %v", err)
	}
	app := App{
		srv:    tprpc.NewServer(),
		db:     db,
		config: config,
	}
	app.Method("session/enter", SessionEnter)
	app.Method("session/leave", SessionLeave)
	app.Method("buffer/enter", BufferEnter)
	app.Method("buffer/leave", BufferLeave)
	app.Method("data/project", GetProjectTime)
	app.Method("data/file", GetMostEditedFile)
	return app
}

type Handler func(db *sqlx.DB, params json.RawMessage) (any, error)

func (app App) Method(method string, handler Handler) {
	srvHandler := func(params json.RawMessage) (any, error) {
		return handler(app.db, params)
	}
	app.srv.Method(method, srvHandler)
}

func (app App) Run() {
	app.srv.ListenAndServe()
}
