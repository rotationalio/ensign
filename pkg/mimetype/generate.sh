#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/ensign/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1beta1 \
    --go_opt=module=github.com/rotationalio/ensign/pkg/mimetype/v1beta1 \
    --go_opt=Mmimetype/v1beta1/mimetype.proto="github.com/rotationalio/ensign/pkg/mimetype/v1beta1;mimetype" \
    mimetype/v1beta1/mimetype.proto
