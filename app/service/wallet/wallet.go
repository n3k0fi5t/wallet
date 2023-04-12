package wallet

import (
	"context"

	mBank "github.com/n3k0fi5t/wallet/app/models/bank"
)

type Service interface {
	// Deposit deposit money to specific user's account
	Deposit(ctx context.Context, accountID string, amount int64) (string, error)

	// Deposit withdraw money from specific user's account
	Withdraw(ctx context.Context, accountID string, amount int64) (string, error)

	// Transfer transfer money from one account to another
	Transfer(ctx context.Context, from, to string, amount int64) (string, error)

	// GetAccount get account information of specific users's account
	GetAccount(ctx context.Context, accountID string) (*mBank.Account, error)
}
