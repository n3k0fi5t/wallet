package wallet

import (
	"context"

	mBank "github.com/n3k0fi5t/wallet/app/models/bank"
	"github.com/n3k0fi5t/wallet/app/repository/bank"
	"github.com/sirupsen/logrus"
)

type dealCategory int

const (
	dealUnkown dealCategory = iota
	dealDeposit
	dealWithdraw
	dealTransfer
)

const (
	PseudoAccount = "c1e395d9-8c00-4124-819a-85b0402900cf"
)

func NewWallet(b bank.Bank) Service {
	return &impl{
		bank: b,
	}
}

type impl struct {
	bank bank.Bank
}

func makeDeal(acc1, acc2 string, amount int64, category dealCategory) *mBank.Dealing {
	switch category {
	case dealDeposit:
		return &mBank.Dealing{
			FromAccountID: PseudoAccount,
			ToAccountID:   acc1,
			Amount:        amount,
		}
	case dealWithdraw:
		return &mBank.Dealing{
			FromAccountID: acc1,
			ToAccountID:   PseudoAccount,
			Amount:        amount,
		}
	case dealTransfer:
		return &mBank.Dealing{
			FromAccountID: acc1,
			ToAccountID:   acc2,
			Amount:        amount,
		}
	default:
		return nil
	}
}

func (im *impl) Deposit(ctx context.Context, accountID string, amount int64) (string, error) {
	deal := makeDeal(accountID, "", amount, dealDeposit)
	tradeID, err := im.bank.Trade(ctx, deal)
	if err != nil {
		logrus.WithField("err", err).Error("bank.Trade failed in Deposit")
		return "", err
	}

	return tradeID, nil
}

func (im *impl) Withdraw(ctx context.Context, accountID string, amount int64) (string, error) {
	deal := makeDeal(accountID, "", amount, dealWithdraw)
	tradeID, err := im.bank.Trade(ctx, deal)
	if err != nil {
		logrus.WithField("err", err).Error("bank.Trade failed in Withdraw")
		return "", err
	}

	return tradeID, nil
}

func (im *impl) Transfer(ctx context.Context, from, to string, amount int64) (string, error) {
	deal := makeDeal(from, to, amount, dealTransfer)
	tradeID, err := im.bank.Trade(ctx, deal)
	if err != nil {
		logrus.WithField("err", err).Error("bank.Trade failed in Transfer")
		return "", err
	}

	return tradeID, nil
}

func (im *impl) GetAccount(ctx context.Context, accountID string) (*mBank.Account, error) {
	var account *mBank.Account
	account, err := im.bank.GetAccount(ctx, accountID)
	if err != nil {
		logrus.WithField("err", err).Error("bank.GetAccount failed in GetAccount")
		return nil, err
	}

	return account, nil
}
