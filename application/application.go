package application

import (
	"crud_web_service/config"
	"crud_web_service/model"
)

// Application tier of 3-layer architecture
type Application struct {
	model  model.Model
	config config.Main
}

// New Application constructor
func New(m model.Model, c config.Main) Application {
	return Application{
		model:  m,
		config: c,
	}
}

// PingDatabase ensures db connection is valid
func (a *Application) PingDatabase() error {
	return a.model.Ping()
}
