//go:build unit

package wallet

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

type StubWallet struct {
	wallets       []Wallet
	walletsByUser []Wallet
	walletsQuery  []Wallet
	createWallet  Wallet
	updateWallet  Wallet
	err           error
}

func (s StubWallet) Wallets() ([]Wallet, error) {
	return s.wallets, s.err
}

func (s StubWallet) WalletsQuery(id string) ([]Wallet, error) {
	return s.walletsQuery, s.err
}

func (s StubWallet) WalletsByUser(id string) ([]Wallet, error) {
	return s.walletsByUser, s.err
}

func (s StubWallet) CreateWallet(wallet Wallet) (Wallet, error) {
	return s.createWallet, s.err
}

func (s StubWallet) UpdateWallet(wallet Wallet, id string) (Wallet, error) {
	return s.updateWallet, s.err
}

func (s StubWallet) DeleteWallet(id string) error {
	return s.err
}

func TestWallet(t *testing.T) {
	t.Run("given unable to get wallets should return 500 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/wallets")

		stubError := StubWallet{err: echo.ErrInternalServerError}
		p := New(stubError)

		p.WalletsHandler(c)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("given wallet type able to getting wallet should return list of wallets", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/wallets")
		c.SetParamNames("wallet_type")
		c.SetParamValues("Savings")

		stubUser := StubWallet{
			walletsQuery: []Wallet{
				{UserName: "John Doe", WalletName: "John Savings", WalletType: "Savings"},
				{UserName: "Jane Doe", WalletName: "Jane Savings", WalletType: "Savings"},
			},
		}
		p := New(stubUser)

		p.WalletsTypeQueryHandler(c)

		wantWalletType := "Savings"
		want := []Wallet{
			{UserName: "John Doe", WalletName: "John Savings", WalletType: wantWalletType},
			{UserName: "Jane Doe", WalletName: "Jane Savings", WalletType: wantWalletType},
		}
		gotJson := rec.Body.Bytes()
		var got []Wallet
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("given user able to getting wallet should return list of wallets", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/wallets/user/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		stubUser := StubWallet{
			walletsByUser: []Wallet{
				{ID: 1, UserID: 1, UserName: "John Doe", WalletType: "Savings"},
				{ID: 2, UserID: 1, UserName: "John Doe", WalletType: "Credit Card"},
				{ID: 3, UserID: 1, UserName: "John Doe", WalletType: "Crypto Wallet"},
			},
		}
		p := New(stubUser)

		p.WalletsByUserHandler(c)

		wantUserName := "John Doe"
		want := []Wallet{
			{ID: 1, UserID: 1, UserName: wantUserName, WalletType: "Savings"},
			{ID: 2, UserID: 1, UserName: wantUserName, WalletType: "Credit Card"},
			{ID: 3, UserID: 1, UserName: wantUserName, WalletType: "Crypto Wallet"},
		}
		gotJson := rec.Body.Bytes()
		var got []Wallet
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
}
