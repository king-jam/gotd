module github.com/king-jam/gotd

go 1.12

require (
	github.com/alexbyk/panicif v1.0.2
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/nlopes/slack v0.5.0
	github.com/pkg/errors v0.8.1 // indirect
	mellium.im/sasl v0.0.0-00010101000000-000000000000 // indirect
)

replace mellium.im/sasl => github.com/mellium/sasl v0.2.1
