package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	mdBank "github.com/n3k0fi5t/wallet/app/models/bank"
	"github.com/n3k0fi5t/wallet/app/service/wallet"
	mockSrv "github.com/n3k0fi5t/wallet/app/service/wallet/mocks"
)

var (
	mockCtx        = context.Background()
	mockAccountID1 = "935f871a-660f-4f19-801e-916c04bb0324"
	mockAccountID2 = "a89b7b78-b9c1-4129-8cff-380bf53f3a49"
	mockAuth1      = "Tim"
	mockAuth2      = "Alex"
	mockTradeID    = "935f871a-660f-4f19-801e-916c04bb0324"
	mockAccount    = &mdBank.Account{
		AccountID: mockAccountID1,
		Balance:   3345678,
	}

	mockAnyCtx              = mock.AnythingOfType("*context.Context")
	mockHandleCtxMiddleware = func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("ctx", mockCtx)
			c.Next()
		}
	}

	anyError = mock.AnythingOfType("error")
)

type testSuite struct {
	suite.Suite

	router  *gin.Engine
	mockSrv *mockSrv.Service
	wsrv    wallet.Service
}

func (s *testSuite) SetupSuite() {
	s.router = gin.Default()
	s.mockSrv = &mockSrv.Service{}
	s.wsrv = s.mockSrv
	handler := NewHandler(s.wsrv)
	rg := s.router.Group("/api/v1")
	rg.Use(mockHandleCtxMiddleware())
	handler.Handle(rg)
}

func (s *testSuite) TearDownSuite() {
	s.mockSrv.AssertExpectations(s.T())
}

func (s *testSuite) SetupTest() {
}

func (s *testSuite) TearDownTest() {
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

// requestHeader return a default header with content type indicating request body type
func requestHeader() http.Header {
	header := http.Header{}

	header.Add("Content-Type", "application/json")
	return header
}

func (s *testSuite) TestDeposit() {
	genPayload := func(d depositParam) []byte {
		b, err := json.Marshal(d)
		s.Require().NoError(err)
		return b
	}

	tests := []struct {
		Desc    string
		Payload []byte
		ExpCode int
		Auth    string
		setup   func()
	}{
		{
			Desc: "normal case",
			setup: func() {
				s.mockSrv.On("Deposit", mockCtx, mockAccountID1, int64(1000)).Return(mockTradeID, nil).Once()
			},
			Payload: genPayload(depositParam{Amount: 1000}),
			Auth:    mockAuth1,
			ExpCode: http.StatusOK,
		},
		{
			Desc: "failed case",
			setup: func() {
				s.mockSrv.On("Deposit", mockCtx, mockAccountID1, int64(1000)).Return("", fmt.Errorf("")).Once()
			},
			Payload: genPayload(depositParam{Amount: 1000}),
			Auth:    mockAuth1,
			ExpCode: http.StatusInternalServerError,
		},
		{
			Desc: "unauthorized case",
			setup: func() {
			},
			Auth:    "",
			ExpCode: http.StatusUnauthorized,
		},
		{
			Desc: "bad param",
			setup: func() {
			},
			Auth:    mockAuth1,
			ExpCode: http.StatusBadRequest,
		},
	}

	for _, t := range tests {
		if t.setup != nil {
			t.setup()
		}

		header := requestHeader()
		header.Set("Authorization", t.Auth)

		req, err := http.NewRequest("POST", "/api/v1/wallet/deposit", bytes.NewBuffer(t.Payload))
		req.Header = header
		s.Require().NoError(err, t.Desc)

		rr := httptest.NewRecorder()
		s.router.ServeHTTP(rr, req)
		s.Require().Equal(t.ExpCode, rr.Code, t.Desc)

	}
}

func (s *testSuite) TestWithdraw() {
	genPayload := func(d withdrawParam) []byte {
		b, err := json.Marshal(d)
		s.Require().NoError(err)
		return b
	}

	tests := []struct {
		Desc    string
		Payload []byte
		ExpCode int
		Auth    string
		setup   func()
	}{
		{
			Desc: "normal case",
			setup: func() {
				s.mockSrv.On("Withdraw", mockCtx, mockAccountID1, int64(1000)).Return(mockTradeID, nil).Once()
			},
			Payload: genPayload(withdrawParam{Amount: 1000}),
			Auth:    mockAuth1,
			ExpCode: http.StatusOK,
		},
		{
			Desc: "failed case",
			setup: func() {
				s.mockSrv.On("Withdraw", mockCtx, mockAccountID1, int64(1000)).Return("", fmt.Errorf("")).Once()
			},
			Payload: genPayload(withdrawParam{Amount: 1000}),
			Auth:    mockAuth1,
			ExpCode: http.StatusInternalServerError,
		},
		{
			Desc: "unauthorized case",
			setup: func() {
			},
			Auth:    "",
			ExpCode: http.StatusUnauthorized,
		},
		{
			Desc: "bad param",
			setup: func() {
			},
			Auth:    mockAuth1,
			ExpCode: http.StatusBadRequest,
		},
	}

	for _, t := range tests {
		if t.setup != nil {
			t.setup()
		}

		header := requestHeader()
		header.Set("Authorization", t.Auth)

		req, err := http.NewRequest("POST", "/api/v1/wallet/withdraw", bytes.NewBuffer(t.Payload))
		req.Header = header
		s.Require().NoError(err, t.Desc)

		rr := httptest.NewRecorder()
		s.router.ServeHTTP(rr, req)
		s.Require().Equal(t.ExpCode, rr.Code, t.Desc)

	}
}

func (s *testSuite) TestTransfer() {
	genPayload := func(d transferParam) []byte {
		b, err := json.Marshal(d)
		s.Require().NoError(err)
		return b
	}

	tests := []struct {
		Desc    string
		Payload []byte
		ExpCode int
		Auth    string
		setup   func()
	}{
		{
			Desc: "normal case",
			setup: func() {
				s.mockSrv.On("Transfer", mockCtx, mockAccountID1, mockAccountID2, int64(1000)).Return(mockTradeID, nil).Once()
			},
			Payload: genPayload(transferParam{ToAccount: mockAccountID2, Amount: 1000}),
			Auth:    mockAuth1,
			ExpCode: http.StatusOK,
		},
		{
			Desc: "failed case",
			setup: func() {
				s.mockSrv.On("Transfer", mockCtx, mockAccountID1, mockAccountID2, int64(1000)).Return("", fmt.Errorf("")).Once()
			},
			Payload: genPayload(transferParam{ToAccount: mockAccountID2, Amount: 1000}),
			Auth:    mockAuth1,
			ExpCode: http.StatusInternalServerError,
		},
		{
			Desc: "unauthorized case",
			setup: func() {
			},
			Auth:    "",
			ExpCode: http.StatusUnauthorized,
		},
		{
			Desc: "bad param",
			setup: func() {
			},
			Auth:    mockAuth1,
			ExpCode: http.StatusBadRequest,
		},
	}

	for _, t := range tests {
		if t.setup != nil {
			t.setup()
		}

		header := requestHeader()
		header.Set("Authorization", t.Auth)

		req, err := http.NewRequest("POST", "/api/v1/wallet/transfer", bytes.NewBuffer(t.Payload))
		req.Header = header
		s.Require().NoError(err, t.Desc)

		rr := httptest.NewRecorder()
		s.router.ServeHTTP(rr, req)
		s.Require().Equal(t.ExpCode, rr.Code, t.Desc)

	}
}

func (s *testSuite) TestGetAccount() {
	tests := []struct {
		Desc       string
		ExpCode    int
		Auth       string
		setup      func()
		ExpAccount accountInfoResp
	}{
		{
			Desc: "normal case",
			setup: func() {
				s.mockSrv.On("GetAccount", mockCtx, mockAccountID1).Return(mockAccount, nil).Once()
			},
			Auth:       mockAuth1,
			ExpCode:    http.StatusOK,
			ExpAccount: accountInfoResp{AccountID: mockAccountID1, Balance: mockAccount.Balance},
		},
		{
			Desc: "failed case",
			setup: func() {
				s.mockSrv.On("GetAccount", mockCtx, mockAccountID1).Return(nil, fmt.Errorf("")).Once()
			},
			Auth:    mockAuth1,
			ExpCode: http.StatusInternalServerError,
		},
		{
			Desc: "unauthorized case",
			setup: func() {
			},
			Auth:    "",
			ExpCode: http.StatusUnauthorized,
		},
	}

	for _, t := range tests {
		if t.setup != nil {
			t.setup()
		}

		header := requestHeader()
		header.Set("Authorization", t.Auth)

		req, err := http.NewRequest("GET", "/api/v1/wallet/account", nil)
		req.Header = header
		s.Require().NoError(err, t.Desc)

		rr := httptest.NewRecorder()
		s.router.ServeHTTP(rr, req)
		s.Require().Equal(t.ExpCode, rr.Code, t.Desc)

		if t.ExpCode == http.StatusOK {
			var resp accountInfoResp
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			s.Require().NoError(err)
			s.Require().Equal(t.ExpAccount, resp, t.Desc)
		}
	}
}
