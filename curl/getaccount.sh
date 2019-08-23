#!/bin/bash

token=$(cat token)

curl -H "Content-Type: application/json" -H "Authorization: BEARER ${token}" -X GET http://localhost:8080/v1/users/theOne
