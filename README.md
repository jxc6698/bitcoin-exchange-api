
## Description
  This repository is bitcoin trade library written by Golang and currently support part of bitmex api. 

  The websocket part is changed from https://github.com/santacruz123/bitmex-go, I have added restart machenism to help recover from disconnected websocket connections.  
  
  restfulApi part use a lot of structures in github.com/BitMEX/api-connectors.git  


## Installation
  go get github.com/jxc6698/bitcoin-exchange-api


## Usage
  You can look at test/bitmex_test.go


## Testing
  add your apikey and secretkey in bitmex_test.go

```
  export BITMEX_API_KEY="api-key"
  export BITMEX_API_SECRET="api-secret"
  go test test/bitmex_test.go
```

***
## LICENCE
MIT LICENCE


***
## Thanks

