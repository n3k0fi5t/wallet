package wallet

import (
	"context"
	"testing"

	mdBank "github.com/n3k0fi5t/wallet/app/models/bank"
	"github.com/n3k0fi5t/wallet/app/repository/bank"
	mockBank "github.com/n3k0fi5t/wallet/app/repository/bank/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	mockCtx        = context.Background()
	mockAccountID1 = "n3k0fi5t"
	mockAccountID2 = "deadbeef"
	mockTradeID    = "935f871a-660f-4f19-801e-916c04bb0324"
	mockDealing    = mdBank.Dealing{}
	mockAccount    = &mdBank.Account{
		AccountID: mockAccountID1,
		Balance:   3345678,
	}

	anyDealing = mock.AnythingOfType("*bank.Dealing")
)

type testSuite struct {
	suite.Suite
	srv   Service
	mBank *mockBank.Bank
}

func (s *testSuite) SetupSuite() {
	s.mBank = &mockBank.Bank{}
	s.srv = NewWallet(s.mBank)
}

func (s *testSuite) TearDownSuite() {
}

func (s *testSuite) SetupTest() {
}

func (s *testSuite) TearDownTest() {
	s.mBank.AssertExpectations(s.T())
}

func (s *testSuite) TestTransfer() {
	tests := []struct {
		Desc       string
		From       string
		To         string
		Amount     int64
		ExpTradeID string
		ExpError   error
		setup      func()
	}{
		{
			Desc:       "normal Path",
			From:       mockAccountID1,
			To:         mockAccountID2,
			Amount:     100,
			ExpTradeID: mockTradeID,
			ExpError:   nil,
			setup: func() {
				s.mBank.On("Trade", mockCtx, anyDealing).Return(mockTradeID, nil).Once()
			},
		},
		{
			Desc:       "bad Path account Not Exist",
			From:       mockAccountID1,
			To:         mockAccountID2,
			Amount:     100,
			ExpTradeID: "",
			ExpError:   bank.ErrAccountNotExist,
			setup: func() {
				s.mBank.On("Trade", mockCtx, anyDealing).Return("", bank.ErrAccountNotExist).Once()
			},
		},
		{
			Desc:       "bad Path, Balance Not Enough",
			From:       mockAccountID1,
			To:         mockAccountID2,
			Amount:     100,
			ExpTradeID: "",
			ExpError:   bank.ErrBalanceNotEnough,
			setup: func() {
				s.mBank.On("Trade", mockCtx, anyDealing).Return("", bank.ErrBalanceNotEnough).Once()
			},
		},
		{
			Desc:       "bad Path, self transfer",
			From:       mockAccountID1,
			To:         mockAccountID2,
			Amount:     100,
			ExpTradeID: "",
			ExpError:   bank.ErrSelfTransfer,
			setup: func() {
				s.mBank.On("Trade", mockCtx, anyDealing).Return("", bank.ErrSelfTransfer).Once()
			},
		},
		{
			Desc:       "bad Path,update Balance fail",
			From:       mockAccountID1,
			To:         mockAccountID2,
			Amount:     100,
			ExpTradeID: "",
			ExpError:   bank.ErrUpdateBalance,
			setup: func() {
				s.mBank.On("Trade", mockCtx, anyDealing).Return("", bank.ErrUpdateBalance).Once()
			},
		},
		{
			Desc:       "bad Path, invalid dealing",
			From:       mockAccountID1,
			To:         mockAccountID2,
			Amount:     100,
			ExpTradeID: "",
			ExpError:   bank.ErrInvalidDealing,
			setup: func() {
				s.mBank.On("Trade", mockCtx, anyDealing).Return("", bank.ErrInvalidDealing).Once()
			},
		},
	}

	for _, test := range tests {
		s.SetupTest()
		if test.setup != nil {
			test.setup()
		}

		tradeID, err := s.srv.Transfer(mockCtx, test.From, test.To, test.Amount)
		s.Require().Equal(test.ExpTradeID, tradeID, test.Desc)
		s.Require().Equal(test.ExpError, err, test.Desc)

		s.TearDownTest()
	}
}

func (s *testSuite) TestWithdraw() {
	tests := []struct {
		Desc       string
		Account    string
		Amount     int64
		ExpTradeID string
		ExpError   error
		setup      func()
	}{
		{
			Desc:       "normal Path",
			Account:    mockAccountID1,
			Amount:     100,
			ExpTradeID: mockTradeID,
			ExpError:   nil,
			setup: func() {
				s.mBank.On("Trade", mockCtx, anyDealing).Return(mockTradeID, nil).Once()
			},
		},
		{
			Desc:       "bad Path",
			Account:    mockAccountID1,
			Amount:     100,
			ExpTradeID: "",
			ExpError:   bank.ErrAccountNotExist,
			setup: func() {
				s.mBank.On("Trade", mockCtx, anyDealing).Return("", bank.ErrAccountNotExist).Once()
			},
		},
	}

	for _, test := range tests {
		s.SetupTest()
		if test.setup != nil {
			test.setup()
		}

		tradeID, err := s.srv.Withdraw(mockCtx, test.Account, test.Amount)
		s.Require().Equal(test.ExpTradeID, tradeID, test.Desc)
		s.Require().Equal(test.ExpError, err, test.Desc)

		s.TearDownTest()
	}
}

func (s *testSuite) TestGetAccount() {
	tests := []struct {
		Desc       string
		Account    string
		expAccount *mdBank.Account
		ExpError   error
		setup      func()
	}{
		{
			Desc:       "normal Path",
			Account:    mockAccountID1,
			expAccount: mockAccount,
			ExpError:   nil,
			setup: func() {
				s.mBank.On("GetAccount", mockCtx, mockAccountID1).Return(mockAccount, nil).Once()
			},
		},
		{
			Desc:       "bad Path",
			Account:    mockAccountID1,
			expAccount: nil,
			ExpError:   bank.ErrAccountNotExist,
			setup: func() {
				s.mBank.On("GetAccount", mockCtx, mockAccountID1).Return((*mdBank.Account)(nil), bank.ErrAccountNotExist).Once()
			},
		},
	}

	for _, test := range tests {
		s.SetupTest()
		if test.setup != nil {
			test.setup()
		}

		account, err := s.srv.GetAccount(mockCtx, test.Account)
		s.Require().Equal(test.expAccount, account, test.Desc)
		s.Require().Equal(test.ExpError, err, test.Desc)

		s.TearDownTest()
	}
}

func (s *testSuite) TestDeposit() {
	tests := []struct {
		Desc       string
		Account    string
		Amount     int64
		ExpTradeID string
		ExpError   error
		setup      func()
	}{
		{
			Desc:       "normal Path",
			Account:    mockAccountID1,
			Amount:     100,
			ExpTradeID: mockTradeID,
			ExpError:   nil,
			setup: func() {
				s.mBank.On("Trade", mockCtx, anyDealing).Return(mockTradeID, nil).Once()
			},
		},
		{
			Desc:       "bad Path",
			Account:    mockAccountID1,
			Amount:     100,
			ExpTradeID: "",
			ExpError:   bank.ErrAccountNotExist,
			setup: func() {
				s.mBank.On("Trade", mockCtx, anyDealing).Return("", bank.ErrAccountNotExist).Once()
			},
		},
	}

	for _, test := range tests {
		s.SetupTest()
		if test.setup != nil {
			test.setup()
		}

		tradeID, err := s.srv.Deposit(mockCtx, test.Account, test.Amount)
		s.Require().Equal(test.ExpTradeID, tradeID, test.Desc)
		s.Require().Equal(test.ExpError, err, test.Desc)

		s.TearDownTest()
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
