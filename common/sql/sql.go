package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Transactx wraps sqlx trasaction in one function and provide error handling
func Transactx(context context.Context, db *sqlx.DB, txFunc func(*sqlx.Tx) error) error {
	var err error

	tx, err := db.BeginTxx(context, nil)
	if err != nil {
		return err
	}

	defer func() {
		// handle panic
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}

		// rollback or commit
		if err != nil {
			if e := tx.Rollback(); e != nil {
				logrus.WithField("err", e).Error("sqlx transaction Rollback() fail")
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			logrus.WithField("err", err).Error("sqlx transaction Commit() fail")
		}
	}()

	err = txFunc(tx)
	return err
}
