#! /bin/sh

# Add User
curl -X POST -F "username=sato kaito" -F "password=123" http://localhost:8080/signup

# Obtain Token 
token=$(curl -X POST -F "username=sato kaito" -F "password=123" http://localhost:8080/login | jq -r '.token')
echo "token: ${token}"

curl -X POST -H "Authorization: ${token}" -F "group_name=new_gruop" http://localhost:8080/auth/add_group

