-- name: CreateUser :one
INSERT INTO users (id, name, email, password, default_group_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;

-- name: Follow :exec
INSERT INTO follow_rels (follower_id, followee_id) VALUES ($1, $2);

-- name: CreateMeigen :one
INSERT INTO meigens (meigen, whom_id, group_id, poet_id) VALUES ($1, $2, $3, $4) RETURNING id;

-- name: CreateGroup :one
INSERT INTO groups (name) VALUES ($1) RETURNING id;

-- name: AddUserToGroup :exec
INSERT INTO user_group_rels (user_id, group_id, permission) VALUES ($1, $2, $3);

-- name: InitDefaultUG :exec
INSERT INTO user_group_rels (user_id, group_id, permission) VALUES ($1, $2, 0xffff);

-- name: CreatePoet :one
INSERT INTO poets (name, group_id) VALUES ($1, $2) RETURNING id;


-- name: DeleteGroup :exec
DELETE FROM groups WHERE id = $1;


-- name: SearchUsers :many
SELECT id, name FROM users WHERE name LIKE $1;

-- name: GetDefaultGroupID :one
SELECT default_group_id FROM users WHERE id = $1;

-- name: GetPoetID :one
SELECT id FROM poets WHERE name = $1 AND group_id = $2;

-- name: CheckPoetExists :one
SELECT count(*) FROM poets WHERE name = $1 AND group_id = $2;

-- name: GetGroupsParticipated :many
SELECT group_id from user_group_rels WHERE user_id = $1;

-- name: CheckUserExistsGroup :one
SELECT count(*) from user_group_rels WHERE user_id = $1 AND group_id = $2;

-- name: CheckGroupExists :one
SELECT count(*) FROM user_group_rels JOIN groups ON user_group_rels.group_id = groups.id
    WHERE user_id = $1 AND groups.name = $2;

-- name: CheckUserExists :one
SELECT count(*) FROM users WHERE id = $1;

-- name: GetUserByName :one
SELECT id FROM users WHERE id = $1;

-- name: Login :one
SELECT * FROM users WHERE id = $1 AND password = $2;

-- name: GetUsernameByID :one
SELECT name FROM users WHERE id = $1;

