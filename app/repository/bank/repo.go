package bank

import (
	"context"
	"fmt"

	mBank "github.com/n3k0fi5t/wallet/app/models/bank"
)

var (
	// ErrAccountNotExist means query account not exist
	ErrAccountNotExist = fmt.Errorf("Account not exist")

	// ErrBalanceNotEnough means account does not have enough money
	ErrBalanceNotEnough = fmt.Errorf("Balance not enough")

	// ErrSelfTransfer means trade want to transfer money to the same account
	ErrSelfTransfer = fmt.Errorf("Self transfer")

	// ErrUpdateBalance
	ErrUpdateBalance = fmt.Errorf("fail to update balance")

	// ErrInvalidDealing
	ErrInvalidDealing = fmt.Errorf("Invalid dealing")
)

type Bank interface {
	// Trade executes dealings
	Trade(ctx context.Context, dealing *mBank.Dealing) (string, error)

	// GetAccount get account Information
	GetAccount(ctx context.Context, accountID string) (*mBank.Account, error)
}
