-- +goose Up
create table users(
    id uuid not null primary key,
    email text not null unique,
    username text not null,
    profile_pic_url text not null,
    password text not null,
    role_id uuid not null references roles(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table users;