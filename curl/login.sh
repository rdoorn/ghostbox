#!/bin/bash

curl -H "Content-Type: application/json" -X POST -d '{"email":"dummy@ixxi.io","password":"Password"}' http://localhost:8080/v1/login -o output
cut -f8 -d\" output > token
rm output
