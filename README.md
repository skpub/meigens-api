# meigens-api (名言API)

名言を管理するAPIです

This is the API you can manage 'meigens(名言)'

# API Documentation

## Authentication Middleware

### `POST` /signup

Query Parameters

* `username`
* `user_id`
* `email`
* `password`

### `POST` /login

Query Parameters

* `user_id`
* `password`

Response &rArr;
```json
{
	"message": "You got an access token.",
	"token": YOUR TOKEN
}
```

## Application

**Authentication required.**

Specify your token `Authentication: [YOUR TOKEN]` in the request header.

### `POST` /auth/search_users

Query Parameters

* `query`

Response &rArr;
```json
{
	"users": [{USER_ID, USERNAME}]
}
```

***

### `GET` /auth/fetch_group_ids

Response &rArr;
```json
{
    "group_ids": [GROUP_ID]
}
```

***


### `POST` /auth/add_group

Query Parameters

* `group_name`

Response &rArr;
```json
{
	"message": "Successfully added the group.",
	"group_id": GROUP_ID
}
```

***

### `POST` /auth/add_meigen

Query Parameters

* `poet`
* `meigen`

Response &rArr;
```json
{
	"message": "Successfully added the meigen.",
	"group_id": MEIGEN_ID
}
```

***

### `POST` /auth/add_meigen_to_group

Query Parameters

* `group_id` (Can be obtained by calling `/auth/fetch_group_ids`)
* `poet`
* `meigen`

Response &rArr;
```json
{
	"message": "Successfully added the meigen.",
	"group_id": MEIGEN_ID
}
```

***

### `POST` /auth/follow

Query Parameters

* `target_id`

***

### `PATCH` /auth/path_user_image

Query Parameters

* `image` (Image File (png, jpg))

Response &rArr;
```json
{
	"message": "Successfully patched the user image.",
}
```

***

### `PATCH` /auth/path_group_image

Query Parameters

* `group_id`
* `image` (Image File (png, jpg))

Response &rArr;
```json
{
	"message": "Successfully patched the group image.",
}
```

***

## ER Diagram
![](DB_ER.png)
