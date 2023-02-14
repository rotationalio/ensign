#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO=$(realpath "$DIR/../..")
DOTENV="$REPO/.env"

GIT_REVISION=$(git rev-parse --short HEAD)
TAG=${GIT_REVISION}
PLATFORM="linux/amd64"

if [ -f $DOTENV ]; then
    set -o allexport
    source $DOTENV
    set +o allexport
fi

OPTIND=1
while getopts :t:p: opt; do
    case $opt in
        t)  TAG=$OPTARG
            ;;
        p)  PLATFORM=$OPTARG
            ;;
        \?) exit 2
            ;;
    esac
done
shift "$((OPTIND-1))"

docker buildx build \
    --platform $PLATFORM \
    -t rotationalio/beacon:$TAG -f $DIR/beacon/Dockerfile \
    --build-arg REACT_APP_TENANT_BASE_URL="https://api.rotational.app/v1/" \
    --build-arg REACT_APP_QUARTERDECK_BASE_URL="https://auth.rotational.app/v1/" \
    --build-arg REACT_APP_ANALYTICS_ID=${REACT_APP_ANALYTICS_ID} \
    --build-arg REACT_APP_VERSION_NUMBER=${TAG} \
    --build-arg REACT_APP_GIT_REVISION=${GIT_REVISION} \
    --build-arg REACT_APP_SENTRY_DSN=${REACT_APP_SENTRY_DSN} \
    --build-arg REACT_APP_SENTRY_ENVIRONMENT=production \
    --build-arg REACT_APP_USE_DASH_LOCALE="false" \
    $REPO

docker tag rotationalio/beacon:$TAG gcr.io/rotationalio-habanero/beacon:$TAG
docker push rotationalio/beacon:$TAG
docker push gcr.io/rotationalio-habanero/beacon:$TAG