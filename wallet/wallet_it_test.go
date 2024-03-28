//go:build integration
// +build integration

package wallet

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type Handler struct {
	store Storer
}

type Storer interface {
	Wallets() ([]Wallet, error)
	WalletsByUser(id string) ([]Wallet, error)
	WalletsQuery(name string) ([]Wallet, error)
	CreateWallet(wallet Wallet) (Wallet, error)
	UpdateWallet(wallet Wallet, id string) (Wallet, error)
	DeleteWallet(id string) error
}

const serverPort = 80

func TestITListWallet(t *testing.T) {

	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgres://root:password@postgres/wallet?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h = New(db)

		e.GET("/wallets/1", h.WalletsByUserHandler)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/wallets/1", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()
	// Assertions
	expected := "[{\"ID\":1,\"UserID\":1,\"UserName\":\"John Doe\",\"WalletName\":\"John Savings\",\"WalletType\":\"Savings\",\"Balance\":1000}]"

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
