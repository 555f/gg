#!/usr/bin/env bash

VERSION=$(gitver show)
PLATFORMS=("windows/amd64" "windows/386" "darwin/arm64" "linux/arm64" "linux/amd64")
PACKAGE_NAME="github.com/555f/gg"

go mod download

for platform in "${PLATFORMS[@]}"
do	
    platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}

    output_name=${GOOS}-${GOARCH}

    if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi

    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-X main.Version=${VERSION}" -o ./build/${output_name} ./cmd/gg
    if [ $? -ne 0 ]; then
   		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi    
done

go-selfupdate ./build ${VERSION}
