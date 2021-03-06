module github.com/mohitkumar/finch/worker

go 1.18

require (
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/mohitkumar/finch v0.0.4
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.47.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/mohitkumar/finch => ../
