#!/bin/bash

curl -H "Content-Type: application/json" -X POST -d '{"firstname":"Firstname","lastname":"Lastname","email":"dummy@ixxi.io","password":"Password"}' http://localhost:8080/v1/users
