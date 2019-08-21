#!/bin/bash

token=$(cat token)

curl -H "Content-Type: application/json" -H "Authorization: BEARER ${token}" -X GET http://localhost:8080/v1/users/username/activate/fb3a85b2-d7fb-4433-88d8-da3ae27cb3fa
