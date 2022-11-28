package raft

//go:generate protoc -I=$GOPATH/src/github.com/rotationalio/ensign/proto --go_out=. --go_opt=module=github.com/rotationalio/ensign/pkg/raft/api --go-grpc_out=. --go-grpc_opt=module=github.com/rotationalio/ensign/pkg/raft/api raft/v1beta1/log.proto raft/v1beta1/raft.proto
