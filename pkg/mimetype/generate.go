package mimetype

//go:generate protoc -I=$GOPATH/src/github.com/rotationalio/ensign/proto --go_out=./v1beta1 --go_opt=module=github.com/rotationalio/ensign/pkg/mimetype/v1beta1 --go_opt=Mmimetype/v1beta1/mimetype.proto=github.com/rotationalio/ensign/pkg/mimetype/v1beta1;mimetype mimetype/v1beta1/mimetype.proto
