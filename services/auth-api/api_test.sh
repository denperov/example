#!/usr/bin/env bash

api_url='http://127.0.0.1:8081'
echo ${api_url}

api_get()  { \
 printf 'request:  %s\n' "$1" 1>&2; \
 result=$(curl -s "${api_url}/$1"); \
 printf 'response: %s\n\n' "$(echo ${result} | jq -c -C '.')" 1>&2; \
 echo ${result}; \
}

api_post() { \
 printf 'request:  %s %s\n' "$1" "$(echo "$2" | jq -c -C '.')" 1>&2; \
 result=$(curl -s "${api_url}/$1" -H 'Content-Type: application/json;charset=UTF-8' --data-binary "$2"); \
 printf 'response: %s\n\n' "$(echo ${result} | jq -c -C '.')" 1>&2; \
 echo ${result}; \
}

result=$(api_post "registration" '{"Login":"alice", "Password":"123"}')
result=$(api_post "registration" '{"Login":"bob", "Password":"123"}')
result=$(api_post "registration" '{"Login":"alice", "Password":"123"}')
result=$(api_post "registration" '{"Login":"alice", "Password":"xxx"}')

result=$(api_post "login" '{"Login":"alice", "Password":"xxx"}')
result=$(api_post "login" '{"Login":"bob", "Password":"xxx"}')

result=$(api_post "login" '{"Login":"alice", "Password":"123"}')
auth_token_1=$(echo ${result} | jq -r '.Result.Token')

result=$(api_post "login" '{"Login":"bob", "Password":"123"}')
auth_token_2=$(echo ${result} | jq -r '.Result.Token')

result=$(api_post "logout" '{"Token": "'${auth_token_1}'"}"')
result=$(api_post "logout" '{"Token": "'${auth_token_2}'"}"')
