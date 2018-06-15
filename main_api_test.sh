#!/usr/bin/env bash

api_url='http://127.0.0.1:8080'
echo ${api_url}

cookie_file=$(mktemp)

pipe=/tmp/pipe_for_curl_out
trap "rm -f ${pipe}; rm -f ${cookie_file}" EXIT

api_get()  { \
 printf '\nrequest: %s\n' "$1" >&2; \
 status=$(curl -b "${cookie_file}" -c "${cookie_file}" -w "%{http_code}" -o ${pipe} -s "${api_url}$1"); \
	response=$(cat ${pipe}); \
	printf 'response: %s %s\n' "${status}" "$(echo ${response} | jq -c -C '.')" >&2; \
}

api_post() { \
 printf '\nrequest:  %s %s\n' "$1" "$(echo "$2" | jq -c -C '.')" >&2; \
 status=$(curl -b "${cookie_file}" -c "${cookie_file}" -w "%{http_code}" -o ${pipe} -s "${api_url}$1" -H 'Content-Type: application/json' -d "$2"); \
	response=$(cat ${pipe}); \
 printf 'response: %s %s\n' "${status}" "$(echo ${response} | jq -c -C '.')" >&2; \
}

response() { \
	echo ${response} | jq -r "${1}"
}

assert_status() { \
	if [ "${status}" -eq "$1" ]
	then
		echo -e '\033[0;32mOK\033[0m'
	else
		echo -e '\033[0;31mUnexpected status: '${status}'\033[0m'
	fi
}

#------------------------------------------------------------------------------

echo Проверка авторизации
api_post "/accounts/logout"
assert_status 401
api_post "/accounts/delete"
assert_status 401
api_get  "/items/list"
assert_status 401
api_post "/items/create"
assert_status 401
api_post "/items/delete"
assert_status 401
api_post "/offers/create"
assert_status 401
api_post "/offers/confirm"
assert_status 401

echo Алиса регистрируется
api_post "/accounts/create" '{"Login":"alice","Password":"123"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/accounts/logout"
assert_status 200

echo Боб регистрируется
api_post "/accounts/create" '{"Login":"bob","Password":"321"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/accounts/logout"
assert_status 200

echo Чак регистрируется
api_post "/accounts/create" '{"Login":"chuck","Password":"1337"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/accounts/logout"
assert_status 200

echo Алиса создаёт три товара
api_post "/accounts/login" '{"Login":"alice","Password":"123"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/items/create" '{"Content":"aaa-aaa-aaa"}'
assert_status 200
item_a=$(response '.ItemID')
api_post "/items/create" '{"Content":"bbb-bbb-bbb"}'
assert_status 200
item_b=$(response '.ItemID')
api_get  "/items/list"
assert_status 200
api_post "/accounts/logout"
assert_status 200

echo Чак пытается удалить товар Алисы
api_post "/accounts/login" '{"Login":"chuck","Password":"1337"}'
assert_status 200
api_post "/items/delete" '{"ItemID":"'${item_a}'"}'
assert_status 400
api_post "/accounts/logout"
assert_status 200

echo Алиса удаляет товар
api_post "/accounts/login" '{"Login":"alice","Password":"123"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/items/delete" '{"ItemID":"'${item_a}'"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/accounts/logout"
assert_status 200

echo Чак пытается создать предложение для самого себя
api_post "/accounts/login" '{"Login":"chuck","Password":"1337"}'
assert_status 200
api_post "/items/create" '{"Content":"ccc-ccc-ccc"}'
assert_status 200
item_c=$(response '.ItemID')
api_post "/offers/create" '{"ItemID":"'${item_c}'","NewOwnerID":"chuck"}'
assert_status 400
api_post "/items/delete" '{"ItemID":"'${item_c}'"}'
assert_status 200
api_post "/accounts/logout"
assert_status 200

echo Чак пытается создать и подтвердить предложение на товар Алисы
api_post "/accounts/login" '{"Login":"chuck","Password":"1337"}'
assert_status 200
api_post "/offers/create" '{"ItemID":"'${item_b}'","NewOwnerID":"bob"}'
assert_status 200
code_b=$(response '.ConfirmationCode')
api_post "/offers/confirm" '{"ConfirmationCode":"'${code_b}'"}'
assert_status 400
api_post "/accounts/logout"
assert_status 200

echo Алиса создаёт предложение для Боба
api_post "/accounts/login" '{"Login":"alice","Password":"123"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/offers/create" '{"ItemID":"'${item_b}'","NewOwnerID":"bob"}'
assert_status 200
code_b=$(response '.ConfirmationCode')
api_get  "/items/list"
assert_status 200
api_post "/accounts/logout"
assert_status 200

echo Чак пытается подтвердить предложение на товар Алисы
api_post "/accounts/login" '{"Login":"chuck","Password":"1337"}'
assert_status 200
api_post "/offers/confirm" '{"ConfirmationCode":"'${code_b}'"}'
assert_status 400
api_post "/accounts/logout"
assert_status 200

echo Боб подтверждает предложение Алисы
api_post "/accounts/login" '{"Login":"bob","Password":"321"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/offers/confirm" '{"ConfirmationCode":"'${code_b}'"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/accounts/logout"
assert_status 200

echo Алиса проверяет свой список товаров
api_post "/accounts/login" '{"Login":"alice","Password":"123"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/accounts/logout"
assert_status 200

echo Боб удаляет товар
api_post "/accounts/login" '{"Login":"bob","Password":"321"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/items/delete" '{"ItemID":"'${item_b}'"}'
assert_status 200
api_get  "/items/list"
assert_status 200
api_post "/accounts/logout"
assert_status 200

echo Очистка
api_post "/accounts/login" '{"Login":"alice","Password":"123"}'
assert_status 200
api_post "/accounts/delete"
assert_status 200

api_post "/accounts/login" '{"Login":"bob","Password":"321"}'
assert_status 200
api_post "/accounts/delete"
assert_status 200

api_post "/accounts/login" '{"Login":"chuck","Password":"1337"}'
assert_status 200
api_post "/accounts/delete"
assert_status 200
