package bank

type Action int32

const (
	Action_UNKNOWN_ACTION Action = 0
	Action_INCREASE       Action = 1
	Action_DECREASE       Action = 2
)

const (
	PseudoAccount = "c1e395d9-8c00-4124-819a-85b0402900cf"
)

type Transaction struct {
	AccountID   string `db:"accountID"`
	Action      Action `db:"action"`
	Amount      int64  `db:"amount"`
	TimestampMs int64  `db:"timeMs"`
	TradeID     string `db:"tradeID"`
}

type Dealing struct {
	FromAccountID string
	ToAccountID   string
	Amount        int64
}

func (d *Dealing) IsValid() bool {
	if d == nil {
		return false
	} else if d.FromAccountID == d.ToAccountID {
		return false
	} else if d.Amount < 0 {
		return false
	}
	return true
}
