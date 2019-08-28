#!/bin/bash -e

PREVIOUS_VERSION=$(git describe --tags --always)
RESULT=$(curl -L https://github.com/rdoorn/ixxi/releases/download/${PREVIOUS_VERSION}/golang.version -o golang.version -w "%{http_code}")
PREVIOUS_GOLANG_VERSION=$(cat golang.version)
CURRENT_GOLANG_VERSION=$(go version)

if [ "${RESULT}" != "200" ]; then
    echo "get golang version returned status: ${RESULT}"
    echo "url: https://github.com/rdoorn/ixxi/releases/download/${PREVIOUS_VERSION}/golang.version"
    exit 0
fi

echo "old: [${PREVIOUS_GOLANG_VERSION}]"
echo "new: [${CURRENT_GOLANG_VERSION}]"

if [ "${CURRENT_GOLANG_VERSION}" == "${PREVIOUS_GOLANG_VERSION}" ]; then
    echo "up to date with latest golang version"
    exit 0
fi

echo "new golang version available, rebuilding"

