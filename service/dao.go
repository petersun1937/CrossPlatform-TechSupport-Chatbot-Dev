package service

import (
	"Tg_chatbot/database"

	"google.golang.org/genproto/protobuf/api"
)

type DAO interface {
	CreateUser() error
	CreatePlayer() error
	CreateTask() error
	GetTask(id string) Task
}

type dao struct {
	database database.Database2
	api      api.Api
}

func NewDAO(database database.Database2, api api.Api) DAO {
	return &dao{
		database: database,
		api:      api,
	}
}

// database access object
func (d *dao) CreateUser() error {
	// save into postgres
	d.database.GetDB().Create(model)
	// d.database.GetPostgresDB().Create(model)
}

func (d *dao) CreatePlayer() error {
	// save into mongodb
	d.database.GetMongoDB().Create(model)
}

func (d *dao) CreateTask() error {
	// save mysql
	d.database.GetMongoMySQLDB().Create(model)
}

func (d *dao) GetTask(id string) Task {
	// remote api server
	d.api.GetTak(id)
}
