#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/ensign/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1beta1 --go-grpc_out=./v1beta1 \
    --go_opt=module=github.com/rotationalio/ensign/pkg/api/v1beta1 \
    --go-grpc_opt=module=github.com/rotationalio/ensign/pkg/api/v1beta1 \
    --go_opt=Mmimetype/v1beta1/mimetype.proto="github.com/rotationalio/ensign/pkg/mimetype/v1beta1;mimetype" \
    --go_opt=Mapi/v1beta1/event.proto="github.com/rotationalio/ensign/pkg/api/v1beta1;api" \
    --go_opt=Mapi/v1beta1/topic.proto="github.com/rotationalio/ensign/pkg/api/v1beta1;api" \
    --go_opt=Mapi/v1beta1/ensign.proto="github.com/rotationalio/ensign/pkg/api/v1beta1;api" \
    --go-grpc_opt=Mmimetype/v1beta1/mimetype.proto="github.com/rotationalio/ensign/pkg/mimetype/v1beta1;mimetype" \
    --go-grpc_opt=Mapi/v1beta1/event.proto="github.com/rotationalio/ensign/pkg/api/v1beta1;api" \
    --go-grpc_opt=Mapi/v1beta1/topic.proto="github.com/rotationalio/ensign/pkg/api/v1beta1;api" \
    --go-grpc_opt=Mapi/v1beta1/ensign.proto="github.com/rotationalio/ensign/pkg/api/v1beta1;api" \
    api/v1beta1/event.proto \
    api/v1beta1/topic.proto \
    api/v1beta1/ensign.proto