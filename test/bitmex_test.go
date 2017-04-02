package main

/**
 *   about subscribe message: see https://www.bitmex.com/app/wsAPI
 *
 */

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"testing"

	bmrestfulapi   "github.com/jxc6698/bitcoin-exchange-api/bitmex/restfulApi"
	bmwebsocket "github.com/jxc6698/bitcoin-exchange-api/bitmex/websocketApi"
	"github.com/jxc6698/bitcoin-exchange-api/bitmex"
	"github.com/jxc6698/bitcoin-exchange-api/utils"
	"time"
	"os"
)


var (
	apikey = "apikey"
	apisecret = "apisecret"
)

func init() {
	apikey = os.Getenv("BITMEX_API_KEY")
	apisecret = os.Getenv("BITMEX_API_SECRET")
}

func subscribeContractsOrder(chOrder chan bmwebsocket.WSOrder) {
	var reset int64 = 1
	ch := make(chan int64)

	for {
		if reset == 1 {
			ws := bmwebsocket.NewWS()
			ws.RegisterReStart(ch)

			err := ws.Connect()
			if err != nil {
				fmt.Println("error: " + err.Error())
			}
			chAuth := ws.Auth(apikey, apisecret)
			fmt.Println(<-chAuth)

			ws.SubOrder(chOrder, []bitmex.Contracts{bitmex.XBTUSD,
				bitmex.XBTM17})
			reset = 0
		}
		break;

		select {
		case  reset = <-ch:
			if 0 == reset {
				fmt.Println("ws will stopped")
				break
			}
		case <-time.After(time.Hour):
		}
	}
}

func Test_websocket_order(t *testing.T) {
	var order bmwebsocket.WSOrder
	chOrder := make(chan bmwebsocket.WSOrder, 100)

	subscribeContractsOrder(chOrder)
	for i:=0;i<3;i++ {
		order = <-chOrder
		t.Log(order)
	}
}


/**
 *  subscribe wallet will only receive one message
 */
func subscribeWallet(chWallet chan bitmex.WSWallet) bitmex.WSWallet {
	ws := bmwebsocket.NewWS()
	err := ws.Connect()
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	chAuth := ws.Auth(apikey, apisecret)

	fmt.Println(<-chAuth)

	_ = ws.SubWallet(chWallet)
	wallet := <- chWallet
	return wallet
}

func Test_websocket_wallet(t *testing.T) {
	chWallet := make(chan bitmex.WSWallet)
	wallet := subscribeWallet(chWallet)

	t.Log(wallet)
}


func subscribeContractsQuotes(chQuote chan bmwebsocket.WSQuote) {
	var reset int64 = 1
	ch := make(chan int64)

	for {
		if reset == 1 {
			ws := bmwebsocket.NewWS()
			ws.RegisterReStart(ch)

			err := ws.Connect()
			if err != nil {
				fmt.Println("error: " + err.Error())
			}

			ws.SubQuote(chQuote, []bitmex.Contracts{bitmex.XBTUSD,
				bitmex.XBTM17})
			reset = 0
		}
		break;

		select {
		case  reset = <-ch:
			if 0 == reset {
				fmt.Println("ws will stopped")
				break
			}
		case <-time.After(time.Hour):
		}
	}
}


func Test_websocket_quote(t *testing.T) {
	var quote bmwebsocket.WSQuote
	chQuote := make(chan bmwebsocket.WSQuote, 100)

	subscribeContractsQuotes(chQuote)
	for i:=0;i<3;i++ {
		quote = <-chQuote
		t.Log(quote)
	}
}



func subscribeTrade(cont1, cont2 bitmex.Contracts, chTrade chan bmwebsocket.WSTrade) {

	var reset int64 = 1
	ch := make(chan int64)

	for i:=0;i<5;i++{
		if 1 == reset {
			ws := bmwebsocket.NewWS()
			ws.RegisterReStart(ch)

			err := ws.Connect()
			if err != nil {
				fmt.Println("error: " + err.Error())
			}
			ws.SubTrade(chTrade, []bitmex.Contracts{cont1, cont2})
			reset = 0
		}
		break

		select {
		case reset = <-ch:
			if 0 == reset {
				fmt.Println("ws has disconnected")
				break
			}
		case <-time.After(time.Hour):
		}
	}
}

func Test_websocket_trade(t *testing.T) {

	chTrade := make(chan bmwebsocket.WSTrade, 100)
	subscribeTrade(bitmex.XBTM17, bitmex.XBTUSD, chTrade)

	var trade bmwebsocket.WSTrade
	for i:=0; i<3 ;i++{
		select {
		case trade = <-chTrade:
			t.Log(trade)
		}
	}
}



var (
	orderapi *bmrestfulapi.OrderApi
	configuration *bitmex.Configuration
	account utils.Platform = utils.Platform{}
)

func SendLimitOrderBuy(orderapi *bmrestfulapi.OrderApi, qua, price float64, symbol string) (*bitmex.Order, *bmrestfulapi.APIResponse, error) {
	order, res, err := orderapi.OrderNew(symbol, bitmex.BUY, 0.0, 0.0, float32(qua), price, 0.0, 0.0, 0.0,
		"", "", 0.0, "", "", bitmex.LIMIT, "", "", "", "")

	return order, res, err
}

func SendLimitOrderSell(orderapi *bmrestfulapi.OrderApi, qua, price float64, symbol string) (*bitmex.Order, *bmrestfulapi.APIResponse, error) {
	order, res, err := orderapi.OrderNew(symbol, bitmex.SELL, 0.0, 0.0, float32(qua), price, 0.0, 0.0, 0.0,
		"", "", 0.0, "", "", bitmex.LIMIT, "", "", "", "")

	return order, res, err
}
func SendMarketOrderBuy(orderapi *bmrestfulapi.OrderApi, qua, price float64, symbol string) (*bitmex.Order, *bmrestfulapi.APIResponse, error) {
	order, res, err := orderapi.OrderNew(symbol, bitmex.BUY, 0.0, 0.0, float32(qua), price, 0.0, 0.0, 0.0,
		"", "", 0.0, "", "", bitmex.MARKET, "", "", "", "")

	return order, res, err
}

func SendMarketOrderSell(orderapi *bmrestfulapi.OrderApi, qua, price float64, symbol string) (*bitmex.Order, *bmrestfulapi.APIResponse, error) {
	order, res, err := orderapi.OrderNew(symbol, bitmex.SELL, 0.0, 0.0, float32(qua), price, 0.0, 0.0, 0.0,
		"", "", 0.0, "", "", bitmex.MARKET, "", "", "", "")

	return order, res, err
}

func CancelOrder(orderID string, clOrdID string, text string) (*bitmex.Order, *bmrestfulapi.APIResponse, error) {
	return orderapi.OrderCancel(orderID, clOrdID, text)
}


func Test_restfulapi_buysell(t *testing.T) {
	var err error
	configuration = bitmex.NewConfiguration( bmrestfulapi.APIClientImpl{})
	orderapi = bmrestfulapi.NewOrderApi(configuration)

	account.Apikey = apikey
	account.Secretkey = apisecret
	//orderapi.Configuration.
	orderapi.Configuration.Host = "https://www.bitmex.com"
	orderapi.Configuration.BasePath = "/api/v1"
	orderapi.Configuration.Account = &account
	orderapi.Configuration.ExpireTime = 5


	_, _, err = SendLimitOrderBuy(orderapi, 1, 1,
		string(bitmex.XBTUSD))
	assert.True(t, err == nil, err.Error())
	_, _, err = SendLimitOrderSell(orderapi, 1, 100000,
		string(bitmex.XBTUSD))
	assert.True(t, err == nil, err.Error())
	_, _, err = SendMarketOrderBuy(orderapi, 1, 0,
		string(bitmex.XBTUSD))
	assert.True(t, err == nil, err.Error())
	_, _, err = SendLimitOrderBuy(orderapi, 1, 0,
		string(bitmex.XBTUSD))
	assert.True(t, err == nil, err.Error())
}


func Test_restfulapi_bulk(t *testing.T) {
	configuration = bitmex.NewConfiguration( bmrestfulapi.APIClientImpl{})
	orderapi = bmrestfulapi.NewOrderApi(configuration)

	account.Apikey = apikey
	account.Secretkey = apisecret
	//orderapi.Configuration.
	orderapi.Configuration.Host = "https://www.bitmex.com"
	orderapi.Configuration.BasePath = "/api/v1"
	orderapi.Configuration.Account = &account
	orderapi.Configuration.ExpireTime = 5

	order1 := orderapi.NewOrder(string(bitmex.XBTUSD), bitmex.BUY, 0.0, 0.0, float32(1), 10, 0.0, 0.0, 0.0,
		"", "", 0.0, "", "", bitmex.LIMIT, "", "", "", "")
	order2 := orderapi.NewOrder(string(bitmex.XBTUSD), bitmex.SELL, 0.0, 0.0, float32(1), 11110, 0.0, 0.0, 0.0,
		"", "", 0.0, "", "", bitmex.LIMIT, "", "", "", "")
	orders := make([]bitmex.Order, 0)
	orders = append(orders, *order1)
	orders = append(orders, *order2)

	ords,_,err := orderapi.OrderNewBulk(orders)
	assert.True(t, err == nil)
	t.Log(orders)
	t.Log("\n")
	t.Log(ords)
}

