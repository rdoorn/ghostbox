#!/bin/bash -e

PREVIOUS_VERSION=$(git describe --tags --always)
PREVIOUS_GOLANG_VERSION=$(curl https://github.com/rdoorn/ixxi/releases/download/${PREVIOUS_VERSION}/golang.version -o -)
CURRENT_GOLANG_VERSION=$(go version)

echo "old: [${PREVIOUS_GOLANG_VERSION}]"
echo "new: [${CURRENT_GOLANG_VERSION}]"

if [ "${CURRENT_GOLANG_VERSION}" != "${PREVIOUS_GOLANG_VERSION}" ]; then
    echo "new golang version available, rebuilding"
else
    echo "up to date with latest golang version"
fi
