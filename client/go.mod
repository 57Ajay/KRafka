module github.com/57ajay/KRafka/client

go 1.25.0

require (
	github.com/57ajay/krafka v0.0.0-20250909140532-b8c9902d5d79
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.9
)

require (
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
)

replace github.com/57ajay/krafka/proto => ../
