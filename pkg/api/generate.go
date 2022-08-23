package api

//go:generate protoc -I=$GOPATH/src/github.com/rotationalio/ensign/proto --go_out=. --go_opt=module=github.com/rotationalio/ensign/pkg/api --go-grpc_out=. --go-grpc_opt=module=github.com/rotationalio/ensign/pkg/api ensign/v1beta1/event.proto ensign/v1beta1/topic.proto ensign/v1beta1/ensign.proto
