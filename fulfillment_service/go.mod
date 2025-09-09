module github.com/57ajay/krafka/fullfillment_service

go 1.25.0

replace github.com/57ajay/krafka => ../

require (
	github.com/57ajay/krafka v0.0.0-00010101000000-000000000000
	github.com/confluentinc/confluent-kafka-go/v2 v2.11.1
	google.golang.org/protobuf v1.36.9
)

require (
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/grpc v1.75.0 // indirect
)
