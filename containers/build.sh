#!/bin/bash
# Builds docker images locally and pushes them to DockerHub and GCR

# Print usage and exit
show_help() {
cat << EOF
Usage: ${0##*/} [-h] [-t TAG] [-p PLATFORM] [clean|build|up|deploy]
Builds docker images and pushes them to DockerHub and GCR.
Flags are as follows (getopt required):

    -h           display this help and exit
    -t TAG       tag the images (default is git revision)
    -p PLATFORM  platform to build the images for (default is linux/amd64)

The docker commands are as follows:

    ${0##*/} clean
    ${0##*/} [-t TAG] [-P PLATFORM] deploy

The clean command clears your docker cache to ensure the build is
successful and the deploy command builds and pushes the images to
DockerHub and to GCR for use in a k8s cluser.

Unless otherwise specified TAG is the git hash and PLATFORM is
linux/amd64 when deploying to ensure the correct images are deployed.

NOTE: realpath is required; you can install it on OS X with

    $ brew install coreutils
EOF
}

export GIT_REVISION=$(git rev-parse --short HEAD)

# Parse command line options with getopt
OPTIND=1
TAG=${GIT_REVISION}
PLATFORM="linux/amd64"

while getopts ht:p: opt; do
    case $opt in
        h)
            show_help
            exit 0
            ;;
        t)  TAG=$OPTARG
            ;;
        p)  PLATFORM=$OPTARG
            ;;
        *)
            show_help >&2
            exit 2
            ;;
    esac
done
shift "$((OPTIND-1))"

# Ensure only zero or one arguments are passed to the script
if [[ $# -gt 1 ]]; then
    show_help >&2
    exit 2
fi

if [[ $# -eq 1 ]]; then
    if [[ $1 == "clean" ]]; then
        docker system prune --all
        exit 0
    elif [[ $1 == "deploy" ]]; then
        echo "deploying ensign images"
    else
        show_help >&2
        exit 2
    fi
fi

# Ask the user for a yes/no response.
ask() {
    local prompt default reply

    if [[ ${2:-} = 'Y' ]]; then
        prompt='Y/n'
        default='Y'
    elif [[ ${2:-} = 'N' ]]; then
        prompt='y/N'
        default='N'
    else
        prompt='y/n'
        default=''
    fi

    while true; do

        # Ask the question (not using "read -p" as it uses stderr not stdout)
        echo -n "$1 [$prompt] "

        # Read the answer (use /dev/tty in case stdin is redirected from somewhere else)
        read -r reply </dev/tty

        # Default?
        if [[ -z $reply ]]; then
            reply=$default
        fi

        # Check if the reply is valid
        case "$reply" in
            Y*|y*) return 0 ;;
            N*|n*) return 1 ;;
        esac

    done
}

# Confirm that we're continuing with the tag
if ! ask "Continue with tag $TAG?" N; then
    exit 1
fi

# Helpful variables
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO=$(realpath "$DIR/..")
DOTENV="$REPO/.env"

# Compute "development" version from latest tag
VERSION="$(git describe --abbrev=0)"
VERSION_MAJOR="${VERSION%%\.*}"
VERSION_MINOR="${VERSION#*.}"
VERSION_MINOR="${VERSION_MINOR%.*}"
VERSION_PATCH="${VERSION##*.}"
VERSION_PATCH=$((VERSION_PATCH+1))

export REACT_APP_VERSION_NUMBER="${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}"
export REACT_APP_GIT_REVISION=$GIT_REVISION

# Load .env file from project root if it exists
if [ -f $DOTENV ]; then
    set -o allexport
    source $DOTENV
    set +o allexport
fi

# Build the primary images
docker buildx build --platform $PLATFORM -t rotationalio/ensign:$TAG -f $DIR/ensign/Dockerfile --build-arg GIT_REVISION=${GIT_REVISION} $REPO
if [ $? -ne 0 ]; then exit 1; fi

docker buildx build --platform $PLATFORM -t rotationalio/tenant:$TAG -f $DIR/tenant/Dockerfile --build-arg GIT_REVISION=${GIT_REVISION} $REPO
if [ $? -ne 0 ]; then exit 1; fi

docker buildx build --platform $PLATFORM -t rotationalio/quarterdeck:$TAG -f $DIR/quarterdeck/Dockerfile --build-arg GIT_REVISION=${GIT_REVISION} $REPO
if [ $? -ne 0 ]; then exit 1; fi

docker buildx build --platform $PLATFORM -t rotationalio/uptime:$TA -f $DIR/uptime/Dockerfile --build-arg GIT_REVISION=${GIT_REVISION} $REPO
if [ $? -ne 0 ]; then exit 1; fi

# Build Beacon
docker buildx build \
    --platform $PLATFORM \
    -t rotationalio/beacon:$TAG -f $DIR/beacon/Dockerfile \
    --build-arg REACT_APP_TENANT_BASE_URL="https://api.rotational.app/v1/" \
    --build-arg REACT_APP_QUARTERDECK_BASE_URL="https://auth.rotational.app/v1/" \
    --build-arg REACT_APP_ANALYTICS_ID=${REACT_APP_ANALYTICS_ID} \
    --build-arg REACT_APP_VERSION_NUMBER=${REACT_APP_VERSION_NUMBER} \
    --build-arg REACT_APP_GIT_REVISION=${GIT_REVISION} \
    --build-arg REACT_APP_SENTRY_DSN=${REACT_APP_SENTRY_DSN} \
    --build-arg REACT_APP_SENTRY_ENVIRONMENT=production \
    --build-arg REACT_APP_USE_DASH_LOCALE="false" \
    $REPO
if [ $? -ne 0 ]; then exit 1; fi

# Push to DockerHub
docker push rotationalio/ensign:$TAG
docker push rotationalio/tenant:$TAG
docker push rotationalio/quarterdeck:$TAG
docker push rotationalio/uptime:$TAG
docker push rotationalio/beacon:$TAG
if [ $? -ne 0 ]; then exit 1; fi

# Retag the images to push to gcr.io
docker tag rotationalio/ensign:$TAG gcr.io/rotationalio-habanero/ensign:$TAG
docker tag rotationalio/tenant:$TAG gcr.io/rotationalio-habanero/tenant:$TAG
docker tag rotationalio/quarterdeck:$TAG gcr.io/rotationalio-habanero/quarterdeck:$TAG
docker tag rotationalio/uptime:$TAG gcr.io/rotationalio-habanero/uptime:$TAG
docker tag rotationalio/beacon:$TAG gcr.io/rotationalio-habanero/beacon:$TAG
if [ $? -ne 0 ]; then exit 1; fi

# Push to GCR
docker push gcr.io/rotationalio-habanero/ensign:$TAG
docker push gcr.io/rotationalio-habanero/tenant:$TAG
docker push gcr.io/rotationalio-habanero/quarterdeck:$TAG
docker push gcr.io/rotationalio-habanero/uptime:$TAG
docker push gcr.io/rotationalio-habanero/beacon:$TAG
if [ $? -ne 0 ]; then exit 1; fi