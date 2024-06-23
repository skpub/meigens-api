-- name: CreateUser :one
INSERT INTO users (id, name, email, password, default_group_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;

-- name: Follow :exec
INSERT INTO follow_rels (follower_id, followee_id) VALUES ($1, $2);

-- name: CreateMeigen :one
INSERT INTO meigens (meigen, whom_id, group_id, poet_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING meigens.id;

-- name: CreateGroup :one
INSERT INTO groups (name) VALUES ($1) RETURNING id;

-- name: AddUserToGroup :exec
INSERT INTO user_group_rels (user_id, group_id, permission) VALUES ($1, $2, $3);

-- name: InitDefaultUG :exec
INSERT INTO user_group_rels (user_id, group_id, permission) VALUES ($1, $2, 0xffff);

-- name: CreatePoet :one
INSERT INTO poets (name, group_id) VALUES ($1, $2) RETURNING id;

-- name: CreateReaction :one
INSERT INTO reactions (meigen_id, user_id, reaction) VALUES ($1, $2, $3) RETURNING id;

-- name: PatchUserImage :one
UPDATE groups SET img = $2 WHERE id = (
    SELECT default_group_id FROM users WHERE users.id = $1)
    RETURNING id;

-- name: PatchGroupImage :exec
UPDATE groups SET img = $2 WHERE id = $1;


-- name: DeleteGroup :exec
DELETE FROM groups WHERE id = $1;


-- name: FetchTL :many
SELECT meigens.meigen, meigens.whom_id, poets.name, meigens.created_at FROM meigens
    JOIN follow_rels ON meigens.whom_id = follow_rels.followee_id OR meigens.whom_id = follow_rels.follower_id
    JOIN groups ON meigens.group_id = groups.id
    JOIN users ON meigens.whom_id = users.id
    JOIN poets ON meigens.poet_id = poets.id
    WHERE (follow_rels.follower_id = $1 OR meigens.whom_id = $1)
        AND users.default_group_id = groups.id
        AND meigens.created_at < $3 ORDER BY meigens.created_at DESC LIMIT $2;

-- name: GetFollowers :many
SELECT follower_id FROM follow_rels WHERE followee_id = $1 ORDER BY follower_id;

-- name: SearchUsers :many
SELECT users.id, users.name, groups.img FROM users JOIN groups ON users.default_group_id = groups.id WHERE users.name LIKE $1;

-- name: GetDefaultGroupID :one
SELECT default_group_id FROM users WHERE id = $1;

-- name: GetPoetID :one
INSERT INTO poets (name, group_id) VALUES ($1, $2) RETURNING id;

-- name: GetPoetIDGroup :one
SELECT id FROM poets WHERE name = $1 AND group_id = $2;

-- name: CheckPoetExists :one
SELECT id FROM poets WHERE name = $1 AND group_id = $2;

-- name: GetGroupsParticipated :many
SELECT group_id from user_group_rels WHERE user_id = $1;

-- name: CheckUserExistsGroup :one
SELECT permission from user_group_rels WHERE user_id = $1 AND group_id = $2;

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

-- name: CheckMeigenExists :one
SELECT count(*) FROM meigens WHERE id = $1;

-- name: CheckMeigenExistsByMeigen :one
SELECT count(*) FROM meigens JOIN poets ON meigens.poet_id = poets.id
    WHERE meigens.meigen = $1 AND meigens.whom_id = $2 AND meigens.group_id = $3 AND poets.name = $4;

-- name: CheckReactionExists :one
SELECT reaction FROM reactions WHERE meigen_id = $1 AND user_id = $2 AND reaction = $3;

-- name: GetUserImg :one
SELECT img FROM groups WHERE id = (SELECT default_group_id FROM users WHERE users.id = $1);
