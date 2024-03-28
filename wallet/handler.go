package wallet

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

type Storer interface {
	Wallets() ([]Wallet, error)
	WalletsByUser(id string) ([]Wallet, error)
	WalletsQuery(name string) ([]Wallet, error)
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

type Err struct {
	Message string `json:"message"`
}

// WalletHandler
//
//	@Summary		Get all wallets
//	@Description	Get all wallets
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Wallet
//	@Router			/api/v1/wallets [get]
//	@Failure		500	{object}	Err
func (h *Handler) WalletHandler(c echo.Context) error {
	wallets, err := h.store.Wallets()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallets)
}

// WalletByUserHandler
//
//	@Summary		Get wallets by UserID
//	@Description	Get wallets by UserID
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Wallet
//	@Router			/api/v1/users/:id/wallets [get]
//	@Failure		500	{object}	Err
func (h *Handler) WalletByUserHandler(c echo.Context) error {
	id := c.Param("id")
	wallets, err := h.store.WalletsByUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, wallets)
}

// WalletTypeQueryHandler
//
//	@Summary		Get wallets by WalletType
//	@Description	Get wallets by WalletType
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Param			q	query string false "name search by wallet_type"
//	@Success		200	{object}	Wallet
//	@Router			/api/v1/wallets/wallet [get]
//	@Failure		500	{object}	Err
func (h *Handler) WalletTypeQueryHandler(c echo.Context) error {
	name := c.QueryParam("wallet_type")
	wallets, err := h.store.WalletsQuery(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, wallets)
}
