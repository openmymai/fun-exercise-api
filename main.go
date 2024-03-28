package main

import (
	"github.com/labstack/echo/v4"
	"github.com/openmymai/fun-exercise-api/postgres"
	"github.com/openmymai/fun-exercise-api/wallet"

	_ "github.com/openmymai/fun-exercise-api/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title			Wallet API
// @version		1.0
// @description	Sophisticated Wallet API
// @host			localhost:1323
func main() {
	p, err := postgres.New()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	handler := wallet.New(p)
	v1 := e.Group("/api/v1")
	{
		v1.GET("/wallets", handler.WalletHandler)
		v1.GET("/wallets/user/:id", handler.WalletByUserHandler)
		v1.GET("/wallets/wallet", handler.WalletTypeQueryHandler)
	}

	e.Logger.Fatal(e.Start(":1323"))
}
