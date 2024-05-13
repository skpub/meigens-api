#! /bin/sh

# Add User
curl -X POST -F "username=sato kaito" -F "password=123" http://localhost:8080/signup

# Obtain Token 
token=$(curl -X POST -F "username=sato kaito" -F "password=123" http://localhost:8080/login | jq -r '.token')
echo "token: ${token}"

curl -X POST -H "Authorization: ${token}" -F "group_name=new_group" http://localhost:8080/auth/add_group

group_id=$(curl -X GET -H "Authorization: ${token}" -F "username=sato kaito" http://localhost:8080/auth/fetch_group_ids | jq -r '.groups[0]')

echo "group_id: ${group_id}"

curl -X POST -H "Authorization: ${token}" -F "group_id=${group_id}" -F "meigen=meigen" http://localhost:8080/auth/add_meigen_to_group
