#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if ! which go > /dev/null; then
	echo "golang needs to be installed"
	exit 1
fi

bin_dir="_output/bin"
mkdir -p ${bin_dir} || true

echo "**************************************************************"
echo "***************** Running Unit Tests *************************"
echo "**************************************************************"

export TF_ACC=1
go test -v ./secret ./backend ./backend/kms

echo "**************************************************************"
echo "***************** Building Linux Source **********************"
echo "**************************************************************"

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build -v -o ${bin_dir}/terraform-provider-secret-${GOOS}-${GOARCH} main.go

echo "**************************************************************"
echo "***************** Building Mac OS X Source *******************"
echo "**************************************************************"

export GOOS=darwin
export GOARCH=amd64
export CGO_ENABLED=0
go build -v -o ${bin_dir}/terraform-provider-secret-${GOOS}-${GOARCH} main.go
