module github.com/alexjorgef/go-bittrex

go 1.16

require (
	github.com/alexjorgef/signalr v0.1.0
	github.com/google/uuid v1.3.0
	github.com/shopspring/decimal v1.3.1
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/stretchr/testify v1.8.0
)

retract (
	v0.6.3 // For retractions only.
	v0.6.2 // Published accidentally.
	v0.6.0 // Published accidentally.
	v0.5.0 // Published accidentally.
	v0.4.0 // Published accidentally.
	v0.3.0 // Published accidentally.
	v0.2.0 // Published accidentally.
)
