package websocketApi

import (
	"encoding/json"
	"strings"
	"sync"
	"time"
	"errors"

	"github.com/apex/log"

	"golang.org/x/net/websocket"
	"github.com/jxc6698/bitcoin-exchange-api/bitmex"
)

const wsURL = "wss://www.bitmex.com/realtime"

//WS - websocket connection object
type WS struct {
	sync.Mutex
	conn       *websocket.Conn
	log        *log.Logger
	nonce      int64
	key        string
	secret     string
	chTrade    map[chan WSTrade][]bitmex.Contracts
	chQuote    map[chan WSQuote][]bitmex.Contracts
	chOrder    map[chan WSOrder][]bitmex.Contracts
	chPosition map[chan WSPosition][]bitmex.Contracts
	chWallet   chan bitmex.WSWallet
	chSucc     map[string][]chan struct{}

	chStart    chan int64
}

//NewWS - creates new websocket object
func NewWS() *WS {
	return &WS{
		nonce:      time.Now().Unix(),
		chTrade:    make(map[chan WSTrade][]bitmex.Contracts, 0),
		chQuote:    make(map[chan WSQuote][]bitmex.Contracts, 0),
		chOrder:    make(map[chan WSOrder][]bitmex.Contracts, 0),
		chPosition: make(map[chan WSPosition][]bitmex.Contracts, 0),
		chSucc:     make(map[string][]chan struct{}, 0),
	}
}

//Connect - connects
func (ws *WS) Connect() error {
	conn, err := websocket.Dial(wsURL, "", "http://localhost/")

	if err != nil {
		return err
	}

	log.Info("Connected")

	ws.conn = conn

	go ws.read()

	return nil
}

//Disconnect - Disconnects from websocket
func (ws *WS) Disconnect() error {
	log.Info("Disconnecting")
	//TODO Close all channels
	return ws.conn.Close()
}

func (ws *WS) read() {
	for {
		var msg string
		websocket.Message.Receive(ws.conn, &msg)

		log.Debugf("Raw: %v", msg)

		switch {
		case strings.HasPrefix(msg, `{"success"`):
			var success wsSuccess
			json.Unmarshal([]byte(msg), &success)
			log.Debugf("Success: %#v", success)

			if channels, found := ws.chSucc[success.Request["op"]]; found {
				for _, ch := range channels {
					select {
					case ch <- struct{}{}:
					default:
					}
				}
			}

			if channels, found := ws.chSucc[success.Request["args"]]; found {
				for _, ch := range channels {
					select {
					case ch <- struct{}{}:
					default:
					}
				}
			}

		case strings.HasPrefix(msg, `{"info"`):
			var info wsInfo
			json.Unmarshal([]byte(msg), &info)
			log.Infof("Info: %v", info)

		case strings.Contains(msg, `{"table"`):
			var table wsData
			json.Unmarshal([]byte(msg), &table)
			log.Debugf("Table: %#v", table)

			switch table.Table {

			case "trade":
				var trades []WSTrade
				json.Unmarshal(table.Data, &trades)

				log.Debugf("Trades: %#v", trades)

				for _, one := range trades {
					ws.trade(one)
				}

			case "quote":
				var quotes []WSQuote
				json.Unmarshal(table.Data, &quotes)

				log.Debugf("Quotes: %#v", quotes)

				for _, one := range quotes {
					ws.quote(one)
				}

			case "order":
				var orders []WSOrder
				json.Unmarshal(table.Data, &orders)

				log.Debugf("Orders: %#v", orders)

				for _, one := range orders {
					ws.order(one)
				}

			case "position":
				var positions []WSPosition
				json.Unmarshal(table.Data, &positions)

				log.Debugf("Positions: %#v", positions)

				for _, one := range positions {
					ws.position(one)
				}

			case "wallet":
				var wallets []bitmex.WSWallet
				json.Unmarshal(table.Data, &wallets)

				for _, one := range wallets {
					ws.wallet(one)
				}

			default:
				log.Infof("%v", msg)
			}

		default:
			log.Infof("%v", msg)
			log.Info("unknown WS message")
			if nil == ws.chStart {
				ws.fatal(errors.New("Unset restart channel"))
			} else {
				log.Info("restart")
				if len(ws.chStart) == 0 {
					ws.chStart <- 1
				}
			}
		}
	}
}

func (ws *WS) sendTrade(ch chan WSTrade, trade WSTrade) {
	select {
	case ch <- trade:
		log.Debugf("Trade sent: %#v - %#v", ch, trade)
	default:
		log.Debugf("Trade channel busy: %#v", ch)
	}
}

func (ws *WS) sendOrder(ch chan WSOrder, order WSOrder) {
	select {
	case ch <- order:
		log.Debugf("Order sent: %#v - %#v", ch, order)
	default:
		log.Debugf("Order channel busy: %#v", ch)
	}
}

func (ws *WS) sendQuote(ch chan WSQuote, quote WSQuote) {
	select {
	case ch <- quote:
		log.Debugf("Quote sent: %#v - %#v", ch, quote)
	default:
		log.Debugf("Quote channel busy: %#v", ch)
	}
}

func (ws *WS) sendPosition(ch chan WSPosition, position WSPosition) {
	select {
	case ch <- position:
		log.Debugf("Position sent: %#v - %#v", ch, position)
	default:
		log.Debugf("Position channel busy: %#v", ch)
	}
}

func (ws *WS) sendWallet(ch chan bitmex.WSWallet, wallet bitmex.WSWallet) {
	select {
	case ch <- wallet:
		log.Debugf("Wallet sent: %#v - %#v", ch, wallet)
	default:
		log.Debugf("Wallet channel busy: %#v", ch)
	}
}

func (ws *WS) trade(trade WSTrade) {
	for ch, symbols := range ws.chTrade {
		if len(symbols) == 0 {
			ws.sendTrade(ch, trade)
			continue
		}

		for _, oneSymbol := range symbols {
			if oneSymbol == bitmex.Contracts(trade.Symbol) {
				ws.sendTrade(ch, trade)
			}
		}
	}
}

func (ws *WS) order(order WSOrder) {
	for ch, symbols := range ws.chOrder {
		if len(symbols) == 0 {
			ws.sendOrder(ch, order)
			continue
		}

		for _, oneSymbol := range symbols {
			if oneSymbol == bitmex.Contracts(order.Symbol) {
				ws.sendOrder(ch, order)
			}
		}
	}
}

func (ws *WS) position(position WSPosition) {
	for ch, symbols := range ws.chPosition {
		if len(symbols) == 0 {
			ws.sendPosition(ch, position)
			continue
		}

		for _, oneSymbol := range symbols {
			if oneSymbol == bitmex.Contracts(position.Symbol) {
				ws.sendPosition(ch, position)
			}
		}
	}
}

func (ws *WS) quote(quote WSQuote) {
	for ch, symbols := range ws.chQuote {
		if len(symbols) == 0 {
			ws.sendQuote(ch, quote)
			continue
		}

		for _, oneSymbol := range symbols {
			if oneSymbol == bitmex.Contracts(quote.Symbol) {
				ws.sendQuote(ch, quote)
			}
		}
	}
}

func (ws *WS) wallet(wallet bitmex.WSWallet) {
	ws.sendWallet(ws.chWallet, wallet)
}

func (ws *WS) send(msg string) {
	defer ws.Unlock()

	log.Debugf("Writing WS: %#v", string(msg))
	ws.Lock()

	if _, err := ws.conn.Write([]byte(msg)); err != nil {
		ws.fatal(err)
	}
}

func (ws *WS) fatal(err error) {
	ws.Disconnect()
	log.Fatalf("%v", err.Error())
}

//SubTrade - subscribes channel to trades
func (ws *WS) SubTrade(ch chan WSTrade, contracts []bitmex.Contracts) {
	ws.Lock()

	if _, ok := ws.chTrade[ch]; !ok {
		ws.chTrade[ch] = contracts
	} else {
		ws.chTrade[ch] = append(ws.chTrade[ch], contracts...)
	}

	ws.Unlock()

	for _, one := range contracts {
		ws.send(`{"op": "subscribe", "args": "trade:` + string(one) + `"}`)
	}
}

//SubQuote - subscribes to quotes
func (ws *WS) SubQuote(ch chan WSQuote, contracts []bitmex.Contracts) {

	ws.Lock()

	if _, ok := ws.chQuote[ch]; !ok {
		ws.chQuote[ch] = contracts
	} else {
		ws.chQuote[ch] = append(ws.chQuote[ch], contracts...)
	}

	ws.Unlock()

	for _, one := range contracts {
		ws.send(`{"op": "subscribe", "args": "quote:` + string(one) + `"}`)
	}
}

func (ws *WS) RegisterReStart(ch chan int64) {
	ws.chStart = ch
}