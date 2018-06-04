#!/usr/bin/env bash

api_url='http://127.0.0.1:8080'
echo ${api_url}

cookie_file=$(mktemp)

pipe=/tmp/pipe_for_curl_out
trap "rm -f ${pipe}; rm -f ${cookie_file}" EXIT

api_get()  { \
 printf 'request: %s\n' "$1" >&2; \
 status=$(curl -b "${cookie_file}" -c "${cookie_file}" -w "%{http_code}" -o ${pipe} -s "${api_url}$1"); \
	response=$(cat ${pipe}); \
	printf 'response: %s %s\n\n' "${status}" "$(echo ${response} | jq -c -C '.')" >&2; \
}

api_post() { \
 printf 'request:  %s %s\n' "$1" "$(echo "$2" | jq -c -C '.')" >&2; \
 status=$(curl -b "${cookie_file}" -c "${cookie_file}" -w "%{http_code}" -o ${pipe} -s "${api_url}$1" -H 'Content-Type: application/json' -d "$2"); \
	response=$(cat ${pipe}); \
 printf 'response: %s %s\n\n' "${status}" "$(echo ${response} | jq -c -C '.')" >&2; \
}

response() { \
	echo ${response} | jq -r "${1}"
}

#------------------------------------------------------------------------------

echo Проверка авторизации
api_post "/accounts/logout"
api_post "/accounts/delete"
api_get  "/items/list"
api_post "/items/create"
api_post "/items/delete"
api_post "/offers/create"
api_post "/offers/confirm"

echo Алиса регистрируется
api_post "/accounts/create" '{"Login":"alice","Password":"123"}'
api_get  "/items/list"
api_post "/accounts/logout"

echo Боб регистрируется
api_post "/accounts/create" '{"Login":"bob","Password":"321"}'
api_get  "/items/list"
api_post "/accounts/logout"

echo Чак регистрируется
api_post "/accounts/create" '{"Login":"chuck","Password":"1337"}'
api_get  "/items/list"
api_post "/accounts/logout"

echo Алиса создаёт три товара
api_post "/accounts/login" '{"Login":"alice","Password":"123"}'
api_get  "/items/list"
api_post "/items/create" '{"Content":"aaa-aaa-aaa"}'
item_a=$(response '.ItemID')
api_post "/items/create" '{"Content":"bbb-bbb-bbb"}'
item_b=$(response '.ItemID')
api_get  "/items/list"
api_post "/accounts/logout"

echo Чак пытается удалить товар Алисы
api_post "/accounts/login" '{"Login":"chuck","Password":"1337"}'
api_post "/items/delete" '{"ItemID":"'${item_a}'"}'
api_post "/accounts/logout"

echo Алиса удаляет товар
api_post "/accounts/login" '{"Login":"alice","Password":"123"}'
api_get  "/items/list"
api_post "/items/delete" '{"ItemID":"'${item_a}'"}'
api_get  "/items/list"
api_post "/accounts/logout"

echo Чак пытается создать и подтвердить предложение на товар Алисы
api_post "/accounts/login" '{"Login":"chuck","Password":"1337"}'
api_post "/offers/create" '{"ItemID":"'${item_b}'","NewOwnerID":"chuck"}'
code_b=$(response '.ConfirmationCode')
api_post "/offers/confirm" '{"ConfirmationCode":"'${code_b}'"}'
api_post "/accounts/logout"

echo Алиса создаёт предложение для Боба
api_post "/accounts/login" '{"Login":"alice","Password":"123"}'
api_get  "/items/list"
api_post "/offers/create" '{"ItemID":"'${item_b}'","NewOwnerID":"bob"}'
code_b=$(response '.ConfirmationCode')
api_get  "/items/list"
api_post "/accounts/logout"

echo Чак пытается подтвердить предложение на товар Алисы
api_post "/accounts/login" '{"Login":"chuck","Password":"1337"}'
api_post "/offers/confirm" '{"ConfirmationCode":"'${code_b}'"}'
api_post "/accounts/logout"

echo Боб подтверждает предложение Алисы
api_post "/accounts/login" '{"Login":"bob","Password":"321"}'
api_get  "/items/list"
api_post "/offers/confirm" '{"ConfirmationCode":"'${code_b}'"}'
api_get  "/items/list"
api_post "/accounts/logout"

echo Алиса проверяет свой список товаров
api_post "/accounts/login" '{"Login":"alice","Password":"123"}'
api_get  "/items/list"
api_post "/accounts/logout"

echo Боб удаляет товар
api_post "/accounts/login" '{"Login":"bob","Password":"321"}'
api_get  "/items/list"
api_post "/items/delete" '{"ItemID":"'${item_b}'"}'
api_get  "/items/list"
api_post "/accounts/logout"

echo Очистка
api_post "/accounts/login" '{"Login":"alice","Password":"123"}'
api_post "/accounts/delete"

api_post "/accounts/login" '{"Login":"bob","Password":"321"}'
api_post "/accounts/delete"

api_post "/accounts/login" '{"Login":"chuck","Password":"1337"}'
api_post "/accounts/delete"
