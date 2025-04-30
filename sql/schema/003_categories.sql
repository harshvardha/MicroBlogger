-- +goose Up
create table categories(
    id uuid not null primary key,
    category text not null unique,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table categories;