module github.com/king-jam/gotd

go 1.12

require (
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/jinzhu/gorm v1.9.8
	github.com/lib/pq v1.1.0
	github.com/lusis/go-slackbot v0.0.0-20180109053408-401027ccfef5 // indirect
	github.com/lusis/slack-test v0.0.0-20190426140909-c40012f20018 // indirect
	github.com/nlopes/slack v0.5.0
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
)

replace mellium.im/sasl => github.com/mellium/sasl v0.2.1
