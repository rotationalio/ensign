package probez

//go:generate protoc -I=$GOPATH/src/github.com/rotationalio/ensign/proto --go_out=. --go_opt=module=github.com/rotationalio/ensign/pkg/utils/probez --go-grpc_out=. --go-grpc_opt=module=github.com/rotationalio/ensign/pkg/utils/probez grpc/health/v1/health.proto
