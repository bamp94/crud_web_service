package model

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"crud_web_service/config"

	"github.com/GuiaBolso/darwin"
	"github.com/gobuffalo/packr"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Model is data tier of 3-layer architecture
type Model struct {
	db *gorm.DB

	// Used for tracing during building sql query.
	// Must be initialized separately for each query.
	logTrace logrus.Fields
}

var (
	errIDIsNotSpecified = errors.New("Идентификатор не задан")
)

// New Model constructor
func NewFromConfig(config config.Database) Model {
	db, err := gorm.Open("postgres", config.ConnURL())
	if err != nil {
		logrus.WithField("connURL", config.ConnURL()).WithError(err).Fatal("can't open connection with a database")
	}
	if err := db.DB().Ping(); err != nil {
		logrus.WithError(err).Fatal("can't ping connection with a database")
	}
	return Model{db: db}
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

func initLogTrace(trace logrus.Fields) logrus.Fields {
	if trace == nil {
		return make(logrus.Fields)
	}
	return trace
}

// Preload is gorm interface func
func (m *Model) Preload(column string, conditions ...interface{}) *Model {
	trace := initLogTrace(m.logTrace)
	trace["preloadColumn-"+column] = column
	trace["preloadConditions-"+column] = conditions
	return &Model{db: m.db.Preload(column, conditions...), logTrace: trace}
}

// Debug is gorm interface func
func (m *Model) Debug() *Model {
	return &Model{db: m.db.Debug(), logTrace: m.logTrace}
}

// Model is gorm interface func
func (m *Model) Model(value interface{}) *Model {
	trace := initLogTrace(m.logTrace)
	trace["modelValueType"] = fmt.Sprintf("%T", value)
	return &Model{db: m.db.Model(value), logTrace: trace}
}

// Select is gorm interface func
func (m *Model) Select(query interface{}, args ...interface{}) *Model {
	trace := initLogTrace(m.logTrace)
	trace["selectQuery"] = query
	trace["selectArgs"] = args
	return &Model{db: m.db.Select(query, args...), logTrace: trace}
}

// Table is gorm interface func
func (m *Model) Table(name string) *Model {
	trace := initLogTrace(m.logTrace)
	trace["tableName"] = name
	return &Model{db: m.db.Table(name), logTrace: trace}
}

// Limit is gorm interface func
func (m *Model) Limit(limit interface{}) *Model {
	trace := initLogTrace(m.logTrace)
	trace["limit"] = limit
	return &Model{db: m.db.Limit(limit), logTrace: trace}
}

func (m *Model) Set(name string, value interface{}) *Model {
	trace := initLogTrace(m.logTrace)
	var i int
	for {
		if _, ok := trace["setName"+strconv.Itoa(i)]; !ok {
			break
		}
		i++
	}
	trace["setName"+strconv.Itoa(i)] = name
	trace["setValue"+strconv.Itoa(i)] = value
	return &Model{db: m.db.Set(name, value), logTrace: trace}
}
