#!/bin/bash

curl -H "Content-Type: application/json" -X POST -d '{"email":"dummy@ixxi.io","password":"Password"}' http://localhost:8080/v1/login -o output
sed -e 's/.*access_token":"//' output | cut -f1 -d\"  > token
rm output
