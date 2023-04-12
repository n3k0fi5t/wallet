package bank

type Account struct {
	ID        int    `db:"id"`
	AccountID string `db:"accountID"`
	Balance   int64  `db:"balance"`
}
