#!/bin/bash

setup_tests() {
  rm data-testing -rf
  command clear
  f1=$(mktemp)
  echo "1
2
3
4" >$f1
  f2=$(mktemp)
  echo "5
6
7
8
9
10
11
12" >$f2
  config=$(mktemp)
  echo 'server:
  host: localhost
  port: 9999
  mode: development

database:
  driver: sqlite3
  max_connections: 25
  max_idle_connections: 10

auth:
  jwt_secret: xxx
  token_expiry: 24h

storage:
  base_path: data-testing
  versions: versions
  thumbnails: thumbs
  logs: logs
  database: tests.db
  max_file_size: 10MB

thumbnails:
  enable_service: true
  max_width: 100  # px
  max_height: 100 # px
' >$config
  set -m
  scripts/run.sh --config $config &
  serv_pid=$!
  set +m
}

cleanup_tests() {
  rm $f1 $f2 $config
  # -$serv_pid is not a mistake it is used to kill proc gp
  kill -2 -$serv_pid
  waitpid $serv_pid
  exit
}

base_url='http://localhost:9999'

setup_tests
trap "cleanup_tests" INT TERM

name='js-testing'
email='test@example.com'
pass='testing-pass'

new_email="updated-$email"
new_pass="updated-$pass"

# PING
sleep 2s
data=$(curl ${base_url}/ping -Ls)
echo $data
# ---

base_url="$base_url/api/v1"
# SIGNUP
sqlite3 data-testing/tests.db "DELETE FROM users WHERE username = '$name';"
if [ -d "data-testing/$name" ]; then
  rm -r data-testing/$name
fi
data=$(curl ${base_url}/register -Ls -XPOST -d \
  "{
    \"username\": \"$name\",
    \"email\": \"$email\",
    \"password\": \"$pass\"
}")
token=$(echo $data | jq -r '.token')
echo SIGNUP\>\> $data
# ---

# LOGIN
data=$(curl ${base_url}/login -Ls -XPOST -d \
  "{
    \"username\": \"$name\",
    \"password\": \"$pass\"
  }")
token=$(echo $data | jq -r '.token')
echo LOGIN\>\> $data
# ---

# GET ME
data=$(curl ${base_url}/me -Ls -H "Authorization: Bearer ${token}")
echo GET ME\>\> $data
# ---

# UPDATE PASS
data=$(curl ${base_url}/me -Ls -XPOST \
  -H "Authorization: Bearer ${token}" \
  -d "{
  \"old_password\": \"$pass\",
  \"new_password\": \"$new_pass\"
}")
echo UPDATE PASS\>\> $data
# ---

# CHECK UPDATE
data=$(curl ${base_url}/login -Ls -XPOST -d \
  "{
    \"username\": \"$name\",
    \"password\": \"$pass\"
  }")
echo UPDATE PASS CHECK 1\>\> $data
# ---
data=$(curl ${base_url}/login -Ls -XPOST -d \
  "{
    \"username\": \"$name\",
    \"password\": \"$new_pass\"
  }")
token=$(echo $data | jq -r '.token')
echo UPDATE PASS CHECK 2\>\> $data
# ---
data=$(curl ${base_url}/me -Ls -H "Authorization: Bearer ${token}")
echo UPDATE PASS CHECK 3\>\> $data
# ---

# UPDATE EMAIL
data=$(curl ${base_url}/me -Ls -XPOST \
  -H "Authorization: Bearer ${token}" \
  -d "{
  \"email\": \"$new_email\"
}")
echo UPDATE EMAIL\>\> $data
# ---

# CHECK UPDATE
data=$(curl ${base_url}/me -Ls -H "Authorization: Bearer ${token}")
me=$data
echo UPDATE EMAIL CHECK 1\>\> $data
# ---
data=$(curl ${base_url}/login -Ls -XPOST -d \
  "{
    \"username\": \"$name\",
    \"password\": \"$new_pass\"
  }")
echo UPDATE EMAIL CHECK 2\>\> $data
# ---

# CREATE TEMP FILE FOR UPLOAD
data=$(
  curl ${base_url}/asset -Ls -XPOST \
    -H "Authorization: Bearer ${token}" \
    -d "{
    \"parent_id\": \"$(echo $me | jq -r '.root_dir')\",
    \"name\": \"test-file\",
    \"size\": $(wc -c <$f1),
    \"is_dir\": false
  }"
)
echo UPLOAD METADATA\>\> $data
upload_url=$(echo $data | jq -r '.upload_url')
echo UPLOAD URL\>\>$upload_url
data=$(
  curl ${base_url}${upload_url} -Ls -XPATCH \
    -H "Authorization: Bearer ${token}" \
    --data-binary @$f1
)
echo UPLOAD FILE\>\> $data
data=$(
  curl ${base_url}${upload_url} -Ls -XGET \
    -H "Authorization: Bearer ${token}"
)
echo DOWNLOAD FILE\>\> $data

# create version 2
data=$(
  curl ${base_url}${upload_url} -Ls -XPUT \
    -H "Authorization: Bearer ${token}" \
    -d "{
    \"size\": $(wc -c <$f2)
  }"
)
echo UPDATE METADATA\>\> $data
data=$(
  curl ${base_url}${upload_url} -Ls -XPATCH \
    -H "Authorization: Bearer ${token}" \
    --data-binary @$f2
)
echo UPLOAD FILE\>\> $data
data=$(
  curl ${base_url}${upload_url} -Ls -XGET \
    -H "Authorization: Bearer ${token}"
)
echo DOWNLOAD FILE\>\> $data

#  delete file and if all versions are deleted
data=$(
  curl ${base_url}${upload_url} -Ls -XDELETE \
    -H "Authorization: Bearer ${token}"
)
echo DELETE FILE\>\> $data

cleanup_tests
