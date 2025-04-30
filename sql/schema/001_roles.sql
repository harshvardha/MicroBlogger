-- +goose Up
create table Roles(
    id uuid not null primary key,
    role_name text not null unique,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table Roles;