#!/bin/bash
# A helper for common docker compose operations to run Ensign locally.

# Print usage and exit
show_help() {
cat << EOF
Usage: ${0##*/} [-h] [-p PROFILE] [clean|build|up]
A helper for common docker compose operations to run Ensign services
locally. Flags are as follows (getopt required):

    -h          display this help and exit
    -p PROFILE  specify he docker compose profile to use

There are two ways to use this script. Run docker compose:

    ${0##*/} up

Build the images, optionally cleaning the docker cache first:

    ${0##*/} clean
    ${0##*/} build

You can also specify a profile to run a subset of services.
For example, to only run the backend services (e.g. ensign,
tenant, and quarterdeck without beacon):

    ${0##*/} -p backend up

NOTE: realpath is required; you can install it on OS X with

    $ brew install coreutils
EOF
}

# Parse command line options with getopt
OPTIND=1
PROFILE="all"

while getopts hp: opt; do
    case $opt in
        h)
            show_help
            exit 0
            ;;
        p)  PROFILE=$OPTARG
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

# Helpful variables
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO=$(realpath "$DIR/..")
DOTENV="$REPO/.env"

# Set environment variables for the build process
export GIT_REVISION=$(git rev-parse --short HEAD)
export REACT_APP_GIT_REVISION=$GIT_REVISION

# Compute "development" version from latest tag
VERSION="$(git describe --abbrev=0)"
VERSION_MAJOR="${VERSION%%\.*}"
VERSION_MINOR="${VERSION#*.}"
VERSION_MINOR="${VERSION_MINOR%.*}"
VERSION_PATCH="${VERSION##*.}"
VERSION_PATCH=$((VERSION_PATCH+1))

export REACT_APP_VERSION_NUMBER="${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}-dev"

# Load .env file from project root if it exists
if [ -f $DOTENV ]; then
    set -o allexport
    source $DOTENV
    set +o allexport
fi



if [[ $# -eq 1 ]]; then
    if [[ $1 == "clean" ]]; then
        docker system prune --all
        exit 0
    elif [[ $1 == "build" ]]; then
        COMPOSE_PROFILES=$PROFILE docker compose -p ensign -f $DIR/docker-compose.yaml build
        exit 0
    elif [[ $1 == "up" ]]; then
        echo "starting docker compose services"
    else
        show_help >&2
        exit 2
    fi
fi

# By default just bring docker compose up
COMPOSE_PROFILES=$PROFILE docker compose -p ensign -f $DIR/docker-compose.yaml up
