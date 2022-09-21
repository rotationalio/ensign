#!/bin/bash
# A helper script for common docker and docker compose operations

# Print usage and exit
show_help() {
cat << EOF
Usage: ${0##*/} [-h] [-t TAG] [-p PLATFORM] [clean|build|up|deploy]
A helper for common docker and docker compose operations to run
Ensign services locally. Flags are as follows (getopt required):

    -h  display this help and exit

The docker compose commands are as follows:

    ${0##*/} build
    ${0##*/} up

These commands build the images and bring the docker compose system
up with the correct configuration and build arguments.

The docker commands are as follows:

    ${0##*/} clean
    ${0##*/} [-t TAG] [-p PLATFORM] deploy

The clean command clears your docker cache to ensure the build is
successful and the deploy command builds and pushes the images to
DockerHub and to GCR for use in a k8s cluser.

Unless otherwise specified TAG is the git hash and PLATFORM is
linux/amd64 when deploying to ensure the correct images are deployed.

NOTE: realpath is required; you can install it on OS X with

    $ brew install coreutils
EOF
}

# Helpful variables
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO=$(realpath "$DIR/..")
DOTENV="$REPO/.env"

# Set environment variables for the build process
export GIT_REVISION=$(git rev-parse --short HEAD)

# Load .env file from project root if it exists
if [ -f $DOTENV ]; then
    set -o allexport
    source $DOTENV
    set +o allexport
fi

# Parse command line options with getopt
OPTIND=1
TAG=${GIT_REVISION}
PLATFORM="linux/amd64"

while getopts htp: opt; do
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
    elif [[ $1 == "build" ]]; then
        docker compose -p ensign -f $DIR/docker-compose.yaml build
        exit 0
    elif [[ $1 == "up" ]]; then
        docker compose -p ensign -f $DIR/docker-compose.yaml up
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

# Build the primary images
docker buildx build --platform $PLATFORM -t rotationalio/ensign:$TAG -f $DIR/ensign/Dockerfile --build-arg GIT_REVISION=${GIT_REVISION} $REPO

# Retag the images to push to gcr.io
docker tag rotationalio/ensign:$TAG gcr.io/rotationalio-habanero/ensign:$TAG

# Push to DockerHub
docker push rotationalio/ensign:$TAG

# Push to GCR
docker push gcr.io/rotationalio-habanero/ensign:$TAG