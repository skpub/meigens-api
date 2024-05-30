#! /bin/bash

# Add User sato
echo "Add User sato"
curl -s -X POST -w '\n' -F "username=sato kaito" -F "user_id=skpub" -F "email=skpub@testmail.com" -F "password=123" http://localhost:8080/signup

# Add User 1
echo -e "\n# Add User kato"
curl -s -X POST -w '\n' -F "username=kato kiyomasa" -F "user_id=kaki" -F "email=kaki@testmail.com" -F "password=1234" http://localhost:8080/signup

# Obtain Token (sato)
echo -e "\n# Obtain Token"
token=$(curl -s -X POST -F "user_id=skpub" -F "password=123" http://localhost:8080/login | jq -r '.token')
echo "token: ${token}"

# Otain Token (kato)
echo -e "\n Obtain Token (kato)"
token_kato=$(curl -s -X POST -F "user_id=kaki" -F "password=1234" http://localhost:8080/login | jq -r '.token')
echo "token (kato): ${token_kato}"

# Add Group
echo -e "\n# Add Group"
group_id=$(curl -s -X POST -H "Authorization: ${token}" -F "group_name=new_group" http://localhost:8080/auth/add_group | jq -r '.group_id')
echo "group_id: ${group_id}"

# Obtain Group ID
echo -e "\n# Obtain Group ID"
group_ids=$(curl -s -X GET -H "Authorization: ${token}" -F "username=sato kaito" http://localhost:8080/auth/fetch_group_ids | jq -r '.group_ids')
echo "group_ids: ${group_ids}"

group_id2=$(echo ${group_ids} | jq -r '.[0]')

# Add Meigen
echo -e "\n# Add Meigen"
meigen_id=$(curl -s -X POST -w '\n' -H "Authorization: ${token}" -F "meigen=meigen" -F "poet=poepoe" http://localhost:8080/auth/add_meigen | jq -r '.meigen_id')

# Add Meigen to the Group
echo -e "\n# Add Meigen to the Group"
curl -s -X POST -w '\n' -H "Authorization: ${token}" -F "group_id=${group_id2}" -F "meigen=meigen" -F "poet=poepoe" http://localhost:8080/auth/add_meigen_to_group

# Search User
echo -e "\n# Search k (%k%)"
found_users=$(curl -s -X POST -H "Authorization: ${token}" -F "query=k" http://localhost:8080/auth/search_users | jq -r '.found_users')
echo "found_users: ${found_users}"

# sato Follows kaki
echo -e "\n# sato Follows kaki"
curl -s -X POST -w '\n' -H "Authorization: ${token}" -F "target_id=kaki" http://localhost:8080/auth/follow

# patch user image (sato)
echo -e "\n # Patch User Image (sato)"
curl -s -X PATCH -w '\n' -H "Authorization: ${token}" -F "image=@./test.jpg" http://localhost:8080/auth/patch_user_image

# patch group image (sato, new_group)
echo -e "\n # Patch Group Image (sato, new_group)"
curl -s -X PATCH -w '\n' -H "Authorization: ${token}" -F "image=@./test.jpg" -F "group_id=${group_id2}" http://localhost:8080/auth/patch_group_image

# create reaction to the meigen (kato -> sato:poepoe)
echo -e "\n # Create Reaction to the Meigen (kato -> sato:poepoe)"
curl -s -X POST -w '\n' -H "Authorization: ${token_kato}" -F "reaction=0" -F "meigen_id=${meigen_id}" http://localhost:8080/auth/reaction

