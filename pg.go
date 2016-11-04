package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"log"
	"strings"
	"strconv"

	xormcore "github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var (
	pgEngineMaster = &PGEngine{}
	pgEngineMux    = sync.RWMutex{}
)

type PGSession struct {
	*xorm.Session
}

type PGEngine struct {
	*xorm.Engine
	dsn    string
	Status bool
}

func PGConfig(engine *PGEngine,
	user, password,	dbName,	host,	port string, ssl bool,
	level xormcore.LogLevel, showSQL bool) error {
	var err error

	dsn := fmt.Sprintf(
		"user=%s dbname=%s host=%s port=%s",
		user,
		dbName,
		host,
		port,
	)
	if password != "" {
		dsn += fmt.Sprintf(
			" password=%s",
			password,
		)
	}

	if ssl {
		dsn += " sslmode=enable"
	} else {
		dsn += " sslmode=disable"
	}

	if engine.Engine, err = xorm.NewEngine("postgres", dsn); err != nil {
		return err
	}

	engine.dsn = dsn
	engine.SetMaxOpenConns(100)
	engine.SetMaxIdleConns(50)
	engine.Logger().SetLevel(level)
	engine.ShowSQL(showSQL)

	if err = engine.Ping(); err != nil {
		engine.Status = false
		return err
	}
	engine.Status = true
	return nil
}

func PGRunCheckMasterStatus(interval time.Duration) {
	OthersRunPeriodicTask(func() time.Duration{
		pgCheckDBEngineStatus(pgEngineMaster)
		return interval
	}, "PGRunCheckMasterStatus", interval*time.Second)
}

func pgCheckDBEngineStatus(engine *PGEngine) {
	err := engine.Ping()
	if err != nil {
		pgEngineMux.RLock()
		if engine.Status {
			pgEngineMux.RUnlock()
			pgEngineMux.Lock()
			engine.Status = false
			pgEngineMux.Unlock()
		} else {
			pgEngineMux.RUnlock()
		}
		log.Printf("pgCheckDBEngineStatus dsn: %s, error: %s\n", engine.dsn, err.Error())
		return
	}

	pgEngineMux.RLock()
	if !engine.Status {
		pgEngineMux.RUnlock()
		pgEngineMux.Lock()
		engine.Status = true
		pgEngineMux.Unlock()
	} else {
		pgEngineMux.RUnlock()
	}
}

func PGNewMasterSession() *PGSession {
	ms := new(PGSession)
	ms.Session = pgEngineMaster.NewSession()

	return ms
}

func PGGetEngine(master, sync bool) (*PGEngine, error) {
	pgEngineMux.RLock()
	defer pgEngineMux.RUnlock()

	if master {
		if pgEngineMaster.Status {
			return pgEngineMaster, nil
		}

		return nil, fmt.Errorf("DB engine not avaliable")
	}

	// if sync {} //TODO
	return nil, fmt.Errorf("DB engine not avaliable")
}

func pgNewMasterAutoCloseSession() *PGSession {
	return pgNewAutoCloseSession(true, false)
}

func pgNewAutoCloseSession(master, sync bool) *PGSession {
	engine, err := PGGetEngine(master, sync)
	if err != nil {
		engine = pgEngineMaster
	}

	ms := new(PGSession)
	ms.Session = engine.NewSession()
	ms.IsAutoClose = true

	return ms
}

type DBModel interface {
	TableName() string
}

func InsertRow(s *PGSession, m DBModel) (err error) {
	if s == nil {
		s = pgNewMasterAutoCloseSession()
	}
	_, err = s.AllCols().InsertOne(m)

	if err != nil && strings.Index(err.Error(), "duplicate key") >= 0 {
		err = errors.New("DB model duplicated")
	}

	return
}

func InsertMultiRows(s *PGSession, m []interface{}) (err error) {
	var ms *PGSession

	if s == nil {
		ms = PGNewMasterSession()
		defer ms.Close()
		if err = ms.Begin(); err != nil {
			return err
		}
	} else {
		ms = s
	}

	_, err = ms.AllCols().Insert(m...)
	if s == nil {
		if err != nil {
			ms.Rollback()
		} else {
			err = ms.Commit()
		}
	}

	if err != nil && strings.Index(err.Error(), "duplicate key") >= 0 {
		err = errors.New("DB model duplicated")
	}

	return
}

type UniqueDBModel interface {
	TableName() string
	UniqueCond() (string, []interface{})
}

func UpdateDBModel(s *PGSession, m UniqueDBModel) (err error) {
	whereStr, whereArgs := m.UniqueCond()
	if s == nil {
		s = pgNewMasterAutoCloseSession()
	}

	_, err = s.AllCols().Where(whereStr, whereArgs...).Update(m)
	if err != nil && strings.Index(err.Error(), "duplicate key") >= 0 {
		err = errors.New("DB model duplicated")
	}

	return
}

func DeleteDBModel(s *PGSession, m UniqueDBModel) (err error) {
	whereStr, whereArgs := m.UniqueCond()

	if s == nil {
		s = pgNewMasterAutoCloseSession()
	}

	_, err = s.Where(whereStr, whereArgs...).Delete(m)
	return
}

func pgGenSequenceValue(sequence string) (int64, error) {
	results, err := pgEngineMaster.Query("select nextval(?) as next", sequence)
	if err != nil {
		return 0, fmt.Errorf("Generate %s sequence error: %s", sequence, err.Error())
	}

	id, err := strconv.ParseInt(string(results[0]["next"]), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Generate %s sequence error: %s", sequence, err.Error())
	}

	return id, nil
}
