-- name: CreateUser :one
INSERT INTO users (id, name, email, password, default_group_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;

-- name: CreateMeigen :one
INSERT INTO meigens (meigen, whom_id, group_id, poet_id) VALUES ($1, $2, $3, $4) RETURNING id;

-- name: CreateGroup :one
INSERT INTO groups (name) VALUES ($1) RETURNING id;

-- name: AddUserToGroup :exec
INSERT INTO user_group_rels (user_id, group_id, permission) VALUES ($1, $2, $3);

-- name: InitDefaultUG :exec
INSERT INTO user_group_rels (user_id, group_id, permission) VALUES ($1, $2, 0xffff);

-- name: GetDefaultGroupID :one
SELECT default_group_id FROM users WHERE id = $1;

-- name: CreatePoet :one
INSERT INTO poets (name, group_id)
    SELECT name, group_id FROM poets
    WHERE NOT EXISTS (
        SELECT * FROM poets where poets.name = $1 AND poets.group_id = $2
    ) LIMIT 1 RETURNING id;

-- name: GetGroupsParticipated :many
SELECT group_id from user_group_rels WHERE user_id = $1;

-- name: UserEXGroup :one
SELECT count(*) from user_group_rels WHERE user_id = $1 AND group_id = $2;

-- name: GroupEX :one
SELECT count(*) FROM user_group_rels JOIN groups ON user_group_rels.group_id = groups.id WHERE user_id = $1 AND groups.name = $2;

-- name: GetUserByName :one
SELECT id FROM users WHERE id = $1;

-- name: Login :one
SELECT * FROM users WHERE id = $1 AND password = $2;

-- name: GetUsernameByID :one
SELECT name FROM users WHERE id = $1;

-- -- name: PoetExGroup :one
-- SELECT count(*) FROM poet_group_rels JOIN poets ON poet_group_rels.poet_id = poets.id
-- WHERE poets.name = $1 AND poet_group_rels.group_id = $2;

-- -- name: PoetEx :one
-- SELECT poet_id FROM poet_user_rels JOIN poets ON poet_user_rels.poet_id = poets.id
-- WHERE poets.name = $1 AND poet_user_rels.user_id = $2;
