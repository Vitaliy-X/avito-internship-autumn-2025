#!/bin/bash

curl -X GET "http://localhost:8080/team/get?team_name=red" \
  -H "Content-Type: application/json"
