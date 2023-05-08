#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/ensign/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

MODULE="github.com/rotationalio/ensign/pkg/ensign/region/v1beta1"
MOD="github.com/rotationalio/ensign/pkg/ensign/region/v1beta1;region"

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1beta1 \
    --go_opt=module="${MODULE}" \
    --go_opt=Mregion/v1beta1/region.proto="${MOD}" \
    region/v1beta1/region.proto
