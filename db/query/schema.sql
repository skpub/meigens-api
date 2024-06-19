CREATE TABLE groups (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(127) NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    img         bytea
);

CREATE TABLE users (
    id                  varchar(127) PRIMARY KEY NOT NULL,
    name                varchar(127) NOT NULL,
    bio                 varchar(511),
    since               TIMESTAMPTZ DEFAULT NOW(),
    email               varchar(127) NOT NULL UNIQUE,
    password            TEXT NOT NULL,
    default_group_id    UUID REFERENCES groups(id) NOT NULL,
    private             BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE user_group_rels (
    user_id     varchar(127) REFERENCES users(id) NOT NULL,
    group_id    UUID REFERENCES groups(id) NOT NULL,
    permission  SMALLINT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, group_id)
);

CREATE INDEX ON user_group_rels(user_id COLLATE "unicode");

CREATE TABLE follow_rels (
    follower_id     varchar(127) REFERENCES users(id) NOT NULL,
    followee_id     varchar(127) REFERENCES users(id) NOT NULL,
    PRIMARY KEY (follower_id, followee_id)
);

CREATE INDEX ON follow_rels(follower_id COLLATE "unicode");

CREATE TABLE poets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(127) NOT NULL,
    group_id    UUID REFERENCES groups(id) NOT NULL
);

CREATE TABLE meigens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meigen      TEXT NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    whom_id     varchar(127) REFERENCES users(id) NOT NULL,
    group_id    UUID REFERENCES groups(id) NOT NULL,
    poet_id     UUID REFERENCES poets(id) NOT NULL
);

CREATE TABLE reactions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meigen_id   UUID REFERENCES meigens(id) NOT NULL,
    user_id     varchar(127) REFERENCES users(id) NOT NULL,
    reaction    INTEGER NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
