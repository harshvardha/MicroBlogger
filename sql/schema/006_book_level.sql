-- +goose Up
create table book_level(
    id uuid not null primary key,
    level text not null unique,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table book_level;