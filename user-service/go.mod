module github.com/nazibul7/go-grpc-microservices/user-service

go 1.25.11

require (
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/golang-migrate/migrate/v4 v4.19.1
	google.golang.org/grpc v1.81.1
)

require google.golang.org/protobuf v1.36.11 // indirect

require (
	github.com/lib/pq v1.10.9 // indirect
	github.com/nazibul7/go-grpc-microservices/proto v0.0.0
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
)

replace github.com/nazibul7/go-grpc-microservices/proto => ../proto
