CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(127) NOT NULL,
    bio         varchar(511),
    since       TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    email       varchar(127) NOT NULL UNIQUE,
    password    TEXT NOT NULL
);

CREATE TABLE groups (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(127) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    owner_id    UUID REFERENCES users(id) NOT NULL
);

CREATE TABLE user_group_rels (
    user_id     UUID REFERENCES users(id) NOT NULL,
    group_id    UUID REFERENCES groups(id) NOT NULL,
    PRIMARY KEY (user_id, group_id)
);

CREATE TABLE poets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(127) NOT NULL
);

CREATE TABLE poet_group_rels (
    poet_id     UUID REFERENCES poets(id) NOT NULL,
    group_id    UUID REFERENCES groups(id) NOT NULL,
    PRIMARY KEY (poet_id, group_id)
);

CREATE TABLE poet_user_rels (
    poet_id     UUID REFERENCES poets(id) NOT NULL,
    user_id     UUID REFERENCES users(id) NOT NULL,
    PRIMARY KEY (poet_id, user_id)
);

CREATE TABLE meigens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meigen      TEXT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    whom_id     UUID REFERENCES users(id) NOT NULL,
    group_id    UUID REFERENCES groups(id),
    poet_id     UUID REFERENCES poets(id) NOT NULL
);
