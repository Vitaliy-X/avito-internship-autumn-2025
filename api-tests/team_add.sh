#!/bin/bash

curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "red",
    "members": [
      { "user_id": "u1", "username": "Vasya", "is_active": true },
      { "user_id": "u2", "username": "Dasha", "is_active": true },
      { "user_id": "u3", "username": "Masha", "is_active": true }
    ]
  }'
echo
