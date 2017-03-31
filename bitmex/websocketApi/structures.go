package websocketApi

import (
	"encoding/json"
	"time"

	"github.com/satori/go.uuid"
	"github.com/jxc6698/bitcoin-exchange-api/bitmex"
)

//WSTrade - trade structure
type WSTrade struct {
	Size            float64   `json:"size"`
	Price           float64   `json:"price"`
	ForeignNotional float64   `json:"foreignNotional"`
	GrossValue      float64   `json:"grossValue"`
	HomeNotional    float64   `json:"homeNotional"`
	Symbol          string    `json:"symbol"`
	TickDirection   string    `json:"tickDirection"`
	Side            string    `json:"side"`
	TradeMatchID    string    `json:"trdMatchID"`
	Timestamp       time.Time `json:"timestamp"`
}

//WSQuote - quote structure
type WSQuote struct {
	Timestamp time.Time `json:"timestamp"`
	Symbol    bitmex.Contracts `json:"symbol"`
	BidPrice  float64   `json:"bidPrice"`
	BidSize   int64     `json:"bidSize"`
	AskPrice  float64   `json:"askPrice"`
	AskSize   int64     `json:"askSize"`
}

//WSPosition - position structure
type WSPosition struct {
	Timestamp        time.Time `json:"timestamp"`
	Symbol           bitmex.Contracts `json:"symbol"`
	Account          int64     `json:"account"`
	CurrentQty       int64     `json:"currentQty"`
	MarkPrice        float64   `json:"markPrice"`
	SimpleQty        float64   `json:"simpleQty"`
	SimplePnl        float64   `json:"simplePnl"`
	LiquidationPrice float64   `json:"liquidationPrice"`
}

//WSOrder - order structure
type WSOrder struct {
	Timestamp             time.Time `json:"timestamp"`
	TransactTime          time.Time `json:"transactTime"`
	OrderID               uuid.UUID `json:"orderID"`
	OrderQty              int64     `json:"orderQty"`
	Account               int64     `json:"account"`
	DisplayQty            int64     `json:"displayQty"`
	LeavesQty             int64     `json:"leavesQty"`
	CumQty                int64     `json:"cumQty"`
	Price                 float64   `json:"price"`
	SimpleOrderQty        float64   `json:"simpleOrderQty"`
	StopPx                float64   `json:"stopPx"`
	PegOffsetValue        float64   `json:"pegOffsetValue"`
	SimpleCumQty          float64   `json:"simpleCumQty"`
	SimpleLeavesQty       float64   `json:"simpleLeavesQty"`
	AvgPx                 float64   `json:"avgPx"`
	Side                  string    `json:"side"`
	ClOrdID               string    `json:"clOrdID"`
	Symbol                bitmex.Contracts `json:"symbol"`
	PegPriceType          string    `json:"pegPriceType"`
	Currency              bitmex.Contracts `json:"currency"`
	SettlCurrency         bitmex.Contracts `json:"settlCurrency"`
	ExecInst              string    `json:"execInst"`
	ContingencyType       string    `json:"contingencyType"`
	ExDestination         string    `json:"exDestination"`
	OrdStatus             string    `json:"ordStatus"`
	Triggered             string    `json:"triggered"`
	WorkingIndicator      string    `json:"workingIndicator"`
	OrdRejReason          string    `json:"ordRejReason"`
	MultiLegReportingType string    `json:"multiLegReportingType"`
	Text                  string    `json:"text"`
}

type wsData struct {
	Table       string            `json:"table"`
	Action      string            `json:"action"`
	Keys        []string          `json:"keys"`
	Attributes  map[string]string `json:"attributes"`
	Types       map[string]string `json:"keys"`
	ForeignKeys map[string]string `json:"foreignKeys"`
	Data        json.RawMessage
}

type wsSuccess struct {
	Success   bool              `json:"success"`
	Subscribe string            `json:"subscribe"`
	Request   map[string]string `json:"request"`
}

type wsInfo struct {
	Info      string    `json:"info"`
	Version   string    `json:"version"`
	Time      time.Time `json:"timestamp"`
	Docs      string    `json:"docs"`
	Heartbeat bool      `json:"heartbeatEnabled"`
}

type wsError struct {
	Error string `json:"error"`
}
