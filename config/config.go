package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

// Main configuration of the project
type Main struct {
	Port     int      `json:"port" binding:"min=1,max=65535"`
	LogLevel string   `json:"logLevel"`
	DB       Database `json:"database"`
}

// Database configuration
type Database struct {
	Host      string `json:"host"     binding:"required"`
	Port      int    `json:"port"     binding:"min=1,max=65535"`
	User      string `json:"user"     binding:"required"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name"     binding:"required"`
	Endpoint  string `json:"endpoint"     binding:"required"`
	Schema    string `json:"schema"     binding:"required"`
	EnableSSL bool   `json:"enableSSL"`
}

// New Main configuration
func New(path string) Main {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.WithError(err).WithField("configPath", path).Fatal("can't read config file in the selected path")
	}
	var config Main
	if err := json.Unmarshal(body, &config); err != nil {
		logrus.WithError(err).Fatal("can't unmarshal config file as a json object")
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal("invalid 'logLevel' parameter in configuration. Available values: ", logrus.AllLevels)
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true) // adds line number to log message

	return config
}

// DbConnURL returns string URL, which may be used for connect to postgres database.
func (c *Database) ConnURL() string {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
	if !c.EnableSSL {
		url += "?sslmode=disable"
	}
	return url
}
