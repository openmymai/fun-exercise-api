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
	CreateWallet(wallet Wallet) (Wallet, error)
	UpdateWallet(wallet Wallet, id string) (Wallet, error)
	DeleteWallet(id string) error
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
func (h *Handler) WalletsHandler(c echo.Context) error {
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
func (h *Handler) WalletsByUserHandler(c echo.Context) error {
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
func (h *Handler) WalletsTypeQueryHandler(c echo.Context) error {
	name := c.QueryParam("wallet_type")
	wallets, err := h.store.WalletsQuery(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, wallets)
}

// CreateWalletHandler
//
//	@Summary		Create wallet
//	@Description	Create wallet
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Wallet
//	@Router			/api/v1/wallets [post]
//	@Failure		500	{object}	Err
func (h *Handler) CreateWalletHandler(c echo.Context) error {
	w := Wallet{}
	err := c.Bind(&w)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	wallet, err := h.store.CreateWallet(w)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, wallet)
}

// UpdateWalletHandler
//
//	@Summary		Update wallet
//	@Description	Update wallet
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Wallet
//	@Router			/api/v1/wallets [put]
//	@Failure		500	{object}	Err
func (h *Handler) UpdateWalletHandler(c echo.Context) error {
	id := c.Param("id")

	wallet := Wallet{}
	err := c.Bind(&wallet)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	updateWallet, err := h.store.UpdateWallet(wallet, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, updateWallet)
}

// DeleteWalletHandler
//
//	@Summary		Delete wallet
//	@Description	Delete wallet
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Wallet
//	@Router			/api/v1/users/:id/wallets [delete]
//	@Failure		500	{object}	Err
func (h *Handler) DeleteWalletHandler(c echo.Context) error {
	id := c.Param("id")

	err := h.store.DeleteWallet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "Delete "+id+" successful")
}
