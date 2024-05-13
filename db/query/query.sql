-- name: CreateUser :exec
INSERT INTO users (name, email, password) VALUES ($1, $2, $3);

-- name: CreateMeigen :exec
INSERT INTO meigens (meigen, whom_id, group_id, poet_id) VALUES ($1, $2, $3, $4);

-- name: CreateGroup :one
INSERT INTO groups (name, owner_id) VALUES ($1, $2) RETURNING id;

-- name: AddUserToGroup :exec
INSERT INTO user_group_rels (user_id, group_id) VALUES ($1, $2);

-- name: CreatePoet :one
INSERT INTO poets (name) VALUES ($1) RETURNING id;

-- name: CreatePoetGroupRel :exec
INSERT INTO poet_group_rels (poet_id, group_id) VALUES ($1, $2);

-- name: GetGroupsParticipated :many
SELECT group_id from user_group_rels WHERE user_id = $1;

-- name: UserEXGroup :one
SELECT count(*) from user_group_rels WHERE user_id = $1 AND group_id = $2;

-- name: GroupEX :one
SELECT count(*) FROM user_group_rels JOIN groups ON user_group_rels.group_id = groups.id WHERE user_id = $1 AND groups.name = $2;

-- name: GetUserByName :one
SELECT name FROM users WHERE name = $1;

-- name: Login :one
SELECT * FROM users WHERE name = $1 AND password = $2;

-- name: GetUsernameByID :one
SELECT name FROM users WHERE id = $1;

-- name: PoetExGroup :one
SELECT count(*) FROM poet_group_rels JOIN poets ON poet_group_rels.poet_id = poets.id
WHERE poets.name = $1 AND poet_group_rels.group_id = $2;
