#!/bin/bash

# Before run: chmod a+x build.sh

# Local Variables
ARGS=("$@")
SERVICE_NAME="${ARGS[0]}"
OS_TYPE="$(go env GOOS)"
ARCH_TYPE="$(go env GOARCH)"
BUILD_PATH="./reweb"
BUILD_DEST="../build/reweb"
BUILD_NAME="reweb"
COPY_TO=".env"
SERVICE_BUILD=$(date '+%Y%m%d%H%M')
SERVICE_SHORT_COMMIT_ID=$(git describe --always)
#SERVICE_COMMIT_ID=$(git rev-parse HEAD)

if [[ $# -eq 0 ]] ; then
    echo "Service file name not found"
    exit 1
fi

cd $BUILD_PATH
GOOS=$OS_TYPE GOARCH=$ARCH_TYPE go build -ldflags="-w -s -X main.ServiceName=$SERVICE_NAME -X main.ServiceVersion=$SERVICE_SHORT_COMMIT_ID -X main.ServiceBuild=$SERVICE_BUILD -X main.ServiceCommitId=$SERVICE_SHORT_COMMIT_ID" -o "$BUILD_NAME"

if [ -f "$BUILD_NAME" ]; then
    echo "$SERVICE_NAME with $SERVICE_SHORT_COMMIT_ID version build successfully"

    chmod +x $BUILD_NAME

    mkdir -p $BUILD_DEST
    mv -f $BUILD_NAME $BUILD_DEST
    echo "$BUILD_NAME file moved into build folder"

    # Via: https://superuser.com/posts/1329136/revisions
    shopt -s dotglob
    cp -vr $COPY_TO $BUILD_DEST
    echo "Additional files/folders moved into build folder"

else
    echo "Something wrong for building service"
fi