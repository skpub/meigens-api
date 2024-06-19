# meigens-api (名言API)

名言を管理するAPIです。

This is the API you can manage 'meigens(名言)'

# Dependency

* sqlc (to generate ORM code)

* Atlas (to migrate)

* Docker (postgres)

TODO: そのうち全部Dockerで管理するようにしたい.

# Migration

`~$ ./migrate.sh`

# API Documentation

## Authorization Middleware

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
	"token": "YOUR TOKEN"
}
```

## Application

** Authorization required.**

Specify your token `Authorization: [YOUR TOKEN]` in the request header.

### `POST` /auth/search_users

Query Parameters

* `query`

Response &rArr;
```json
{
	"users": ["USER_ID", "USERNAME"]
}
```

***

### `GET` /auth/fetch_group_ids

Response &rArr;
```json
{
    "group_ids": ["GROUP_ID"]
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
	"group_id": "GROUP_ID"
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
	"meigen_id": "MEIGEN_ID"
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
	"meigen_id": "MEIGEN_ID"
}
```

***

### `POST` /auth/follow

Query Parameters

* `target_id`

***

### `PATCH` /auth/patch_user_image

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

### `POST` /auth/reaction

Query Parameters

* `meigen_id`
* `reaction` (int32 (enum))

***
Response &rArr;
```json
{
	"message": "Successfully added the reaction.",
	"reaction_id": REACTION_ID,
}
```

### `GET` /auth/fetch_tl

* `before` (unixtime)

Response &rArr;
```json
{
	"contents": ["MEIGEN", "WHOM(user) ID", "POET NAME"]
}
```

### `GET` /auth/socket => Upgrade to WebSocket
名言をポストした際、ログイン中のフォロワーに名言が送られます。
こんな感じでサーバがクライアントに何かを送信する系の操作を実現するためには双方向通信が必要であり、
名言をポストしたりフェッチしたりする際にはこのソケット(WebSocket)が使われることになります。
(自動更新する名言のTLみたいなのを実現するためにWebSocketを使っているわけです。)

When you post a meigen, it will be sent to your logged-in followers.
To perform this operation, bidirectional communication is required.
Therefore, this endpoint has been upgraded to WebSocket,
and you need to use this socket to communicate with the meigens-api when posting or fetching a meigen.


#### client to server (or server to client) data format:

```csv
[0-1], JSON DATA
```

instructions:

- 0 => Post the meigen as your posts.

payload will be

```json
{
	"meigen": "a ramen is wind.",
	"poet": "sato",
}
```

- 1 => Post the meigen as specified group's post.

payload will be

```json
{
	"meigen": "a ramen is wind.",
	"poet": "sato",
	"group_id": "GROUP_ID"
}
```

***


## ER Diagram
![](DB_ER.png)
