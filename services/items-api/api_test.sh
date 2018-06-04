#!/usr/bin/env bash

api_url='http://127.0.0.1:8082'
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

result=$(api_get "list?OwnerID=alice")

result=$(api_post "create" '{"OwnerID":"alice", "Content":"111111"}')
item_id_1=$(echo ${result} | jq -r '.Result.ItemID')

result=$(api_get "list?OwnerID=alice")

result=$(api_post "create" '{"OwnerID":"alice", "Content":"222222"}')
item_id_2=$(echo ${result} | jq -r '.Result.ItemID')

result=$(api_get "list?OwnerID=alice")

result=$(api_post "transfer" '{"OwnerID":"alice", "ItemID":"'${item_id_2}'", "NewOwnerID":"bob"}')

result=$(api_post "transfer" '{"OwnerID":"alice", "ItemID":"'${item_id_2}'", "NewOwnerID":"bob"}')

result=$(api_post "transfer" '{"OwnerID":"bob", "ItemID":"'${item_id_1}'", "NewOwnerID":"alice"}')

result=$(api_get "list?OwnerID=alice")

result=$(api_get "list?OwnerID=bob")

result=$(api_post "delete" '{"OwnerID": "alice", "ItemID":"'${item_id_1}'"}')

result=$(api_get "list?OwnerID=alice")

result=$(api_post "delete" '{"OwnerID": "bob", "ItemID":"'${item_id_2}'"}')

result=$(api_get "list?OwnerID=bob")
