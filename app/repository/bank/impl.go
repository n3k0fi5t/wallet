package bank

import (
	"context"

	"github.com/jmoiron/sqlx"
	mBank "github.com/n3k0fi5t/wallet/app/models/bank"
	"github.com/n3k0fi5t/wallet/app/util"
	"github.com/n3k0fi5t/wallet/common/sql"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
)

const (
	queryAccount         = "SELECT id, accountID, balance FROM account WHERE accountID = ?"
	queryBalance         = "SELECT balance FROM account WHERE accountID = ?"
	updateBalance        = "UPDATE account SET balance = balance + ? WHERE accountID = ?"
	insertTransactionLog = "INSERT INTO TransactionLog (accountID, action, amount, timestampMs, tradeID) VALUES (?, ?, ?, ?, ?)"
)

var (
	timeNowMs = util.TimeNowMs
)

func NewBank(db *sqlx.DB) Bank {
	return &impl{
		db: db,
	}
}

type impl struct {
	db           *sqlx.DB
	singleflight singleflight.Group
}

func createTradingLog(dealing *mBank.Dealing, tradeID string, timestamp int64) (debit, credit *mBank.Transaction) {
	debit = &mBank.Transaction{
		AccountID:   dealing.FromAccountID,
		Action:      mBank.Action_DECREASE,
		Amount:      dealing.Amount,
		TimestampMs: timestamp,
		TradeID:     tradeID,
	}

	credit = &mBank.Transaction{
		AccountID:   dealing.ToAccountID,
		Action:      mBank.Action_INCREASE,
		Amount:      dealing.Amount,
		TimestampMs: timestamp,
		TradeID:     tradeID,
	}
	return debit, credit
}

func (im *impl) checkAccountBalance(ctx context.Context, tx *sqlx.Tx, accountID string, amount int64) (bool, error) {
	var balance int64
	if err := tx.Get(&balance, queryBalance, accountID); err != nil {
		return false, err
	}

	if balance < amount {
		return false, nil
	}
	return true, nil
}

func (im *impl) updateBalance(ctx context.Context, tx *sqlx.Tx, accountID string, amount int64) error {
	res, err := tx.ExecContext(ctx, updateBalance, amount, accountID)
	if err != nil {
		logrus.WithField("err", err).Error("ExecContext failed")
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		logrus.WithField("err", err).Error("RowsAffected failed")
		return err
	}

	if affected != 1 {
		logrus.WithField("affected", affected).Error("unexpected affected rows")
		return ErrUpdateBalance
	}
	return nil
}

func (im *impl) logTrading(ctx context.Context, tx *sqlx.Tx, dealing *mBank.Dealing, tradeID string, timestampMs int64) error {
	debit, credit := createTradingLog(dealing, tradeID, timestampMs)
	if _, err := tx.ExecContext(ctx, insertTransactionLog, debit.AccountID, debit.Action, debit.Amount, debit.TimestampMs, debit.TradeID); err != nil {
		logrus.WithField("err", err).Error("ExecContext failed")
		return err
	}

	if _, err := tx.ExecContext(ctx, insertTransactionLog, credit.AccountID, credit.Action, credit.Amount, credit.TimestampMs, credit.TradeID); err != nil {
		logrus.WithField("err", err).Error("ExecContext failed")
		return err
	}

	return nil
}

func (im *impl) trade(ctx context.Context, tx *sqlx.Tx, dealing *mBank.Dealing) (string, error) {
	nowMs := timeNowMs()
	tradeID, err := util.GetUUIDv4()
	if err != nil {
		logrus.WithField("err", err).Error("GetUUIDv4 failed in Bank.trade")
		return "", err
	}

	if ok, err := im.checkAccountBalance(ctx, tx, dealing.FromAccountID, dealing.Amount); err != nil {
		logrus.WithField("err", err).Error("checkAccountBalance failed in Bank.trade")
		return "", err
	} else if !ok {
		return "", ErrBalanceNotEnough
	}

	// doing transfer
	if err = im.updateBalance(ctx, tx, dealing.FromAccountID, -1*dealing.Amount); err != nil {
		logrus.WithField("err", err).Error("updateBalance failed in Bank.trade")
		return "", err
	}
	if err = im.updateBalance(ctx, tx, dealing.ToAccountID, dealing.Amount); err != nil {
		logrus.WithField("err", err).Error("updateBalance failed in Bank.trade")
		return "", err
	}

	// write transaction log (double entries)
	if err := im.logTrading(ctx, tx, dealing, tradeID, nowMs); err != nil {
		logrus.WithField("err", err).Error("logTransaction failed in Bank.trade")
		return "", err
	}

	return tradeID, nil
}

func (im *impl) Trade(ctx context.Context, dealing *mBank.Dealing) (string, error) {
	// do dealing check
	if !dealing.IsValid() {
		return "", ErrInvalidDealing
	}

	tradeID := ""
	if err := sql.Transactx(ctx, im.db, func(tx *sqlx.Tx) error {
		tID, err := im.trade(ctx, tx, dealing)
		tradeID = tID
		return err
	}); err != nil {
		return "", err
	}

	return tradeID, nil
}

func (im *impl) getAccount(ctx context.Context, tx *sqlx.Tx, accountID string) (*mBank.Account, error) {
	accounts := []*mBank.Account{}
	if err := tx.SelectContext(ctx, &accounts, queryAccount, accountID); err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, ErrAccountNotExist
	}

	return accounts[0], nil
}

func (im *impl) GetAccount(ctx context.Context, accountID string) (*mBank.Account, error) {
	// use singleflight to avoid spike query for the same accountID
	val, err, _ := im.singleflight.Do(accountID, func() (interface{}, error) {
		var acc *mBank.Account
		if err := sql.Transactx(ctx, im.db, func(tx *sqlx.Tx) error {
			res, e := im.getAccount(ctx, tx, accountID)
			if e != nil {
				return e
			}

			acc = res
			return nil
		}); err != nil {
			return nil, err
		}
		return acc, nil
	})

	// singleflight.Do only return function error. In this case, type conversion failure happends while function return nil
	res, ok := val.(*mBank.Account)
	if !ok {
		return nil, err
	}

	return res, nil
}
