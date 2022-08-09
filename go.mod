module github.com/alexjorgef/go-bittrex

go 1.17

require (
	github.com/alexjorgef/signalr v0.1.0
	github.com/google/uuid v1.3.0
	github.com/shopspring/decimal v1.3.1
	github.com/stretchr/testify v1.8.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
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
