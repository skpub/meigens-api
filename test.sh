#! /bin/bash

# Add User
curl -s -X POST -F "username=sato kaito" -F "user_id=skpub" -F "email=skpub@testmail.com" -F "password=123" http://localhost:8080/signup

# Obtain Token 
token=$(curl -s -X POST -F "user_id=skpub" -F "password=123" http://localhost:8080/login | jq -r '.token')
echo "token: ${token}"

# Add Group
group_id=$(curl -s -X POST -H "Authorization: ${token}" -F "group_name=new_group" http://localhost:8080/auth/add_group | jq -r '.group_id')
echo "group_id: ${group_id}"

# Obtain Group ID
group_ids=$(curl -s -X GET -H "Authorization: ${token}" -F "username=sato kaito" http://localhost:8080/auth/fetch_group_ids | jq -r '.group_ids')
echo "group_ids: ${group_ids}"

group_id2=$(echo ${group_ids} | jq -r '.[0]')

curl -s -X POST -H "Authorization: ${token}" -F "group_id=${group_id2}" -F "meigen=meigen" -F "poet=poepoe" http://localhost:8080/auth/add_meigen_to_group
