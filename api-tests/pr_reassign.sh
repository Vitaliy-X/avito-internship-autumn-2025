#!/bin/bash

curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-2",
    "old_user_id": "u2"
  }'
