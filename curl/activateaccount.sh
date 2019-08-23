#!/bin/bash

token=$(cat token)

curl -H "Content-Type: application/json" -H "Authorization: BEARER ${token}" -X GET http://localhost:8080/v1/users/theOne/activate/52557f8f-bbd3-40b2-984e-5b90e54f4e05 -o output
sed -e 's/.*access_token":"//' output | cut -f1 -d\"  > token
rm output
