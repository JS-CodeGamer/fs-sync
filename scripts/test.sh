#!/bin/bash

command clear

base_url='http://localhost:8080'

# create 10 dummy users

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

# SIGNUP
sqlite3 data/prod.db "DELETE FROM users WHERE username = '$name';"
if [ -d "data/$name" ]; then
  rm -r data/$name
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
f1=$(mktemp)
echo file-1 contents:
echo "1
2
3
4" | tee $f1
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

f2=$(mktemp)
echo file-2 contents:
echo "5
6
7
8
9
10
11
12" | tee $f2
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

# # TEST ADMIN STUFF
# sqlite3 db.sqlite3 "UPDATE users SET role = 'ADMIN' WHERE email = '$email';"
# data=$(curl ${base_url}/admin/users -Ls \
#   -H "Authorization: Bearer ${token}")
# echo ADMIN GET ALL USERS\>\> $data
# # ---
