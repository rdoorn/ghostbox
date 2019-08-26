#!/bin/bash


if [ -z "$1" ]; then
    echo "specify token"
    exit 1
fi
token=$(cat token)

curl -H "Content-Type: application/json" -H "Authorization: BEARER ${token}" -X GET http://localhost:8080/v1/users/theOne/activate/$1 -o output
sed -e 's/.*access_token":"//' output | cut -f1 -d\"  > token
rm output
