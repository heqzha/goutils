package db

import (
	"fmt"
	"strconv"
	"time"

	xormcore "github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	ccc "github.com/heqzha/goutils/concurrency"
	_ "github.com/lib/pq"
)

type PGEngine struct {
	*xorm.Engine
	dsn string
}

type PGMasterEngine struct {
	*PGEngine
}

func (pg *PGMasterEngine) Init(user, password, dbName, host, port string, ssl bool, level xormcore.LogLevel, showSQL bool) error {
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
	pg = &PGMasterEngine{}
	engine, err := xorm.NewEngine("postgres", dsn)
	if err != nil {
		return err
	}
	pg.Engine = engine
	pg.dsn = dsn
	pg.SetMaxOpenConns(100)
	pg.SetMaxIdleConns(50)
	pg.Logger().SetLevel(level)
	pg.ShowSQL(showSQL)

	if err := pg.Ping(); err != nil {
		return err
	}
	return nil
}

func (pg *PGMasterEngine) SyncModels(models ...interface{}) error {
	return pg.Sync2(models...)
}

func (pg *PGMasterEngine) EnableEngineStatusChecker(t time.Duration) {
	ccc.TaskRunPeriodic(func() time.Duration {
		if err := pg.Ping(); err != nil {
			fmt.Println("Cannot connect to %s: %s", pg.dsn, err.Error())
		}
		return t
	}, "PGRunCheckMasterStatus", time.Second)
}

func (pg *PGMasterEngine) GetAutoCloseSession() *xorm.Session {
	s := pg.NewSession()
	s.IsAutoClose = true
	return s
}

func (pg *PGMasterEngine) NewSequence(sequence string) (int64, error) {
	results, err := pg.Query("select nextval(?) as next", sequence)
	if err != nil {
		return 0, fmt.Errorf("Generate sequence %s error: %s", sequence, err.Error())
	} else if results == nil || len(results) == 0 {
		return 0, fmt.Errorf("Generate sequence %s error: results is empty", sequence)
	}

	result, err := strconv.ParseInt(string(results[0]["next"]), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Generate sequence %s error: %s", sequence, err.Error())
	}

	return result, nil
}
