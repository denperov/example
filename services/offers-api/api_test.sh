#!/usr/bin/env bash

api_url='http://127.0.0.1:8083'
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

result=$(api_post "transfer/create" '{"SenderID":"alice", "RecipientID":"bob", "ItemID":"XXXXXXXXXXXXXXXX"}')
conf_code_1=$(echo ${result} | jq -r '.Result.ConfirmationCode')

result=$(api_post "transfer/create" '{"SenderID":"bob", "RecipientID":"alice", "ItemID":"YYYYYYYYYYYYYYYY"}')
conf_code_2=$(echo ${result} | jq -r '.Result.ConfirmationCode')

result=$(api_post "transfer/confirm" '{"RecipientID":"bob", "ConfirmationCode":"'${conf_code_1}'"}')
