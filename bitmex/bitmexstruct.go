package bitmex

import (
	"time"
)

const (
	BUY     = "Buy"
	SELL    = "Sell"
)

const (
	// Market, Limit, Stop, StopLimit, MarketIfTouched, LimitIfTouched, MarketWithLeftOverAsLimit, Pegged
	LIMIT   = "Limit"
	MARKET  = "Market"
	STOP    = "Stop"
	STOPLIMIT  = "StopLimit"
	MARKETIFTOUCHED  = "MarketIfTouched"
	LIMITIFTOUCHED  = "LimitIfTouched"
	MARKETWITHLEFTOVERASLIMIT  = "MarketWithLeftOverAsLimit"
	PEGGED  = "Pegged"
)


type ErrorError struct {
	Message string `json:"message,omitempty"`
	Name    string `json:"name,omitempty"`
}

type ErrorMessage struct {
	Errormsg ErrorError `json:"error,omitempty"`
}

type WSWallet struct {

	Account float32 `json:"account,omitempty"`

	Currency string `json:"currency,omitempty"`

	PrevDeposited float32 `json:"prevDeposited,omitempty"`

	PrevWithdrawn float32 `json:"prevWithdrawn,omitempty"`

	PrevAmount float32 `json:"prevAmount,omitempty"`

	PrevTimestamp time.Time `json:"prevTimestamp,omitempty"`

	DeltaDeposited float32 `json:"deltaDeposited,omitempty"`

	DeltaWithdrawn float32 `json:"deltaWithdrawn,omitempty"`

	DeltaAmount float32 `json:"deltaAmount,omitempty"`

	Deposited float32 `json:"deposited,omitempty"`

	Withdrawn float32 `json:"withdrawn,omitempty"`

	Amount float32 `json:"amount,omitempty"`

	PendingCredit float32 `json:"pendingCredit,omitempty"`

	PendingDebit float32 `json:"pendingDebit,omitempty"`

	ConfirmedDebit float32 `json:"confirmedDebit,omitempty"`

	Timestamp time.Time `json:"timestamp,omitempty"`

	Addr string `json:"addr,omitempty"`

	WithdrawalLock []interface{} `json:"withdrawalLock,omitempty"`
}

type AbstructAPI struct {
	Configuration *Configuration
}