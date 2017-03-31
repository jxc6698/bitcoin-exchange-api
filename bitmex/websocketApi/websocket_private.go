package websocketApi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/jxc6698/bitcoin-exchange-api/bitmex"
)

//Auth - authentication
func (ws *WS) Auth(key, secret string) chan struct{} {
	ws.key = key
	ws.secret = secret

	nonce := ws.Nonce()

	req := fmt.Sprintf("GET/realtime%d", nonce)
	signature := ws.sign(req)

	msg := fmt.Sprintf(
		`{"op": "authKey", "args": ["%s", %d, "%s"]}`,
		key, nonce, signature,
	)

	ch := make(chan struct{})
	ws.Lock()
	ws.chSucc["authKey"] = append(ws.chSucc["authKey"], ch)
	ws.Unlock()

	ws.send(msg)

	return ch
}

func (ws *WS) sign(payload string) string {
	sig := hmac.New(sha256.New, []byte(ws.secret))
	sig.Write([]byte(payload))
	return hex.EncodeToString(sig.Sum(nil))
}

//Nonce - gets next nonce
func (ws *WS) Nonce() int64 {
	ws.nonce++
	return ws.nonce
}

//SubOrder - subscribe to order events
func (ws *WS) SubOrder(ch chan WSOrder, contracts []bitmex.Contracts) chan struct{} {
	ws.Lock()

	if _, ok := ws.chOrder[ch]; !ok {
		ws.chOrder[ch] = contracts
	} else {
		ws.chOrder[ch] = append(ws.chOrder[ch], contracts...)
	}

	ws.Unlock()

	return ws.subPrivate("order")
}

//SubPosition - subscribe to position chage events
func (ws *WS) SubPosition(ch chan WSPosition, contracts []bitmex.Contracts) chan struct{} {
	ws.Lock()

	if _, ok := ws.chPosition[ch]; !ok {
		ws.chPosition[ch] = contracts
	} else {
		ws.chPosition[ch] = append(ws.chPosition[ch], contracts...)
	}

	ws.Unlock()

	return ws.subPrivate("position")
}

func (ws *WS) subPrivate(topic string) chan struct{} {

	ch := make(chan struct{})
	ws.Lock()
	ws.chSucc[topic] = append(ws.chSucc[topic], ch)
	ws.Unlock()

	ws.send(`{"op": "subscribe", "args": "` + topic + `"}`)

	return ch

}

func (ws *WS) SubWallet(ch chan bitmex.WSWallet) chan struct{} {
	ws.Lock()

	ws.chWallet = ch

	ws.Unlock()

	return ws.subPrivate("wallet")
}

func (ws *WS) SubWalletAgain() chan struct{} {
	return ws.subPrivate("wallet")
}