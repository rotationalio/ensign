#!/bin/bash

PROTOS="${GOPATH}/src/github.com/rotationalio/ensign/proto"

if [[ ! -d $PROTOS ]]; then
    echo "cannot find ${PROTOS}"
    exit 1
fi

MODULE="github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
APIMOD="github.com/rotationalio/ensign/pkg/ensign/api/v1beta1;api"
MMEMOD="github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1;mimetype"
REGMOD="github.com/rotationalio/ensign/pkg/ensign/region/v1beta1;region"

# Generate the protocol buffers
protoc -I=${PROTOS} \
    --go_out=./v1beta1 --go-grpc_out=./v1beta1 \
    --go_opt=module=${MODULE} \
    --go-grpc_opt=module=${MODULE} \
    --go_opt=Mmimetype/v1beta1/mimetype.proto="${MMEMOD}" \
    --go_opt=Mregion/v1beta1/region.proto="${REGMOD}" \
    --go_opt=Mapi/v1beta1/event.proto="${APIMOD}" \
    --go_opt=Mapi/v1beta1/topic.proto="${APIMOD}" \
    --go_opt=Mapi/v1beta1/ensign.proto="${APIMOD}" \
    --go_opt=Mapi/v1beta1/groups.proto="${APIMOD}" \
    --go_opt=Mapi/v1beta1/query.proto="${APIMOD}" \
    --go-grpc_opt=Mmimetype/v1beta1/mimetype.proto="${MMEMOD}" \
    --go-grpc_opt=Mregion/v1beta1/region.proto="${REGMOD}" \
    --go-grpc_opt=Mapi/v1beta1/event.proto="${APIMOD}" \
    --go-grpc_opt=Mapi/v1beta1/topic.proto="${APIMOD}" \
    --go-grpc_opt=Mapi/v1beta1/ensign.proto="${APIMOD}" \
    --go-grpc_opt=Mapi/v1beta1/groups.proto="${APIMOD}" \
    --go-grpc_opt=Mapi/v1beta1/query.proto="${APIMOD}" \
    api/v1beta1/event.proto \
    api/v1beta1/topic.proto \
    api/v1beta1/ensign.proto \
    api/v1beta1/groups.proto \
    api/v1beta1/query.proto