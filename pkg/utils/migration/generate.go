package migration

//go:generate protoc -I=$GOPATH/src/github.com/rotationalio/ensign/proto --go_out=. --go_opt=module=github.com/rotationalio/ensign/pkg/utils/migration migration/migration.proto
