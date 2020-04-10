package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"crud_web_service/common"
	"crud_web_service/config"

	"github.com/GuiaBolso/darwin"
	"github.com/gobuffalo/packr"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Model is data tier of 3-layer architecture
type Model struct {
	db     *gorm.DB
	config config.Database
}

// QueryParameters contains query parameters from endpoint
type QueryParameters struct {
	Objects []Object
	Method  string
	Body    string
}

// Object contains object from db table
type Object struct {
	Name string
	ID   string
}

// New Model constructor
func NewFromConfig(config config.Database) Model {
	db, err := gorm.Open("postgres", config.ConnURL())
	if err != nil {
		logrus.WithField("connURL", config.ConnURL()).WithError(err).Fatal("can't open connection with a database")
	}
	if err := db.DB().Ping(); err != nil {
		logrus.WithError(err).Fatal("can't ping connection with a database")
	}
	return Model{db: db, config: config}
}

// CheckMigrations validates database condition
func (m *Model) CheckMigrations() error {
	driver := darwin.NewGenericDriver(m.db.DB(), darwin.PostgresDialect{})
	d := darwin.New(driver, m.getMigrations(), nil)
	if err := d.Validate(); err != nil {
		return err
	}
	migrationInfo, err := d.Info()
	if err != nil {
		return err
	}
	for _, i := range migrationInfo {
		if i.Status == darwin.Applied {
			continue
		}
		return fmt.Errorf("found not applied migration: %s", i.Migration.Description)
	}
	return nil
}

// Migrate applies all migrations to connected database
func (m *Model) Migrate() {
	driver := darwin.NewGenericDriver(m.db.DB(), darwin.PostgresDialect{})
	d := darwin.New(driver, m.getMigrations(), nil)
	if err := d.Migrate(); err != nil {
		logrus.WithError(err).Error("can't migrate")
	}
}

// getMigrations provides migrations in darwin format
func (m *Model) getMigrations() []darwin.Migration {
	// migrationBox is used for embedding the migrations into the binary
	box := packr.NewBox("../etc/migrations")
	var migrations []darwin.Migration
	arr := box.List()
	sort.Strings(arr)
	for i, fileName := range arr {
		if !(strings.HasSuffix(fileName, ".sql") || strings.HasSuffix(fileName, ".SQL")) {
			logrus.Warnf("found file %s with unexpected type, skipping", fileName)
			continue
		}

		migration, err := box.FindString(fileName)
		if err != nil {
			logrus.WithError(err).Error("internal error of packr library")
		}
		migrations = append(migrations, darwin.Migration{
			Version:     float64(i + 1),
			Description: fileName,
			Script:      migration,
		})
	}
	return migrations
}

// Ping connection with database
func (m *Model) Ping() error {
	return m.db.DB().Ping()
}

// ExecQueryByParams executes query to function from database by parameters
func (m *Model) ExecQueryByParams(params QueryParameters) (interface{}, error) {
	var allObjectsNames, allObjectsIDs []string
	for _, object := range params.Objects {
		allObjectsNames = append(allObjectsNames, object.Name)
		allObjectsIDs = append(allObjectsIDs, object.ID)
	}
	var query string
	if params.Body != "" {
		query = "SELECT * FROM " + m.config.Schema + "." + strings.Join(allObjectsNames, "_") + "_" +
			params.Method + "(" + strings.Join(allObjectsIDs, ", ") + ", '" + params.Body + "')"
	}
	if params.Body == "" {
		query = "SELECT * FROM " + m.config.Schema + "." + strings.Join(allObjectsNames, "_") + "_" +
			params.Method + "(" + strings.Join(allObjectsIDs, ", ") + ")"
	}
	if params.Method == common.MethodPost {
		query = "SELECT * FROM " + m.config.Schema + "." + strings.Join(allObjectsNames, "_") + "_" +
			params.Method + "('" + params.Body + "')"
	}
	rows, err := m.db.Raw(query).Rows()
	if err != nil {
		logrus.WithError(err).WithField("params", params).Error("Can't get rows by query")
		return nil, common.ErrInternal
	}

	// Get the column names from the query
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		logrus.WithError(err).WithField("params", params).Error("Can't get get columns")
		return nil, common.ErrInternal
	}

	colNum := len(columns)
	var results []interface{}
	for rows.Next() {
		r := make([]interface{}, colNum)
		for i := range r {
			r[i] = &r[i]
		}

		err = rows.Scan(r...)
		if err != nil {
			logrus.WithError(err).WithField("params", params).Error("Can't scan rows")
			return nil, common.ErrInternal
		}

		var row interface{}
		for i := range r {
			if err := json.Unmarshal(r[i].([]byte), &row); err != nil {
				logrus.WithError(err).WithField("params", params).Error("Can't unmarshal row json")
				return nil, common.ErrInternal
			}
			results = append(results, row)
		}
	}
	if len(results) == 1 {
		return results[0], nil
	}
	return results, nil
}
