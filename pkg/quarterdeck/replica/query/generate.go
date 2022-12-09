package query

//go:generate protoc -I=$GOPATH/src/github.com/rotationalio/ensign/proto --go_out=. --go_opt=module=github.com/rotationalio/ensign/pkg/quarterdeck/replica/query quarterdeck/query/v1beta1/query.proto
