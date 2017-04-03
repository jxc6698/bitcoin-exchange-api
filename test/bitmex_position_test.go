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


var (
	orderapi *bmrestfulapi.OrderApi
	positionapi *bmrestfulapi.PositionApi

	configuration *bitmex.Configuration
	account utils.Platform = utils.Platform{}
)


func Test_restfulapi_position(t *testing.T) {
	var (
		po *bmrestfulapi.Position
		err error
	)

	configuration = bitmex.NewConfiguration( bmrestfulapi.APIClientImpl{})
	positionapi = bmrestfulapi.NewPositionApi(configuration)

	account.Apikey = apikey
	account.Secretkey = apisecret
	//orderapi.Configuration.
	configuration.Host = "https://www.bitmex.com"
	configuration.BasePath = "/api/v1"
	configuration.Account = &account
	configuration.ExpireTime = 5

	po, _, err = positionapi.PositionUpdateLeverage(bitmex.XBTUSD, 2)
	if nil != err {
		fmt.Println(err)
	}
	fmt.Println(po)

}

