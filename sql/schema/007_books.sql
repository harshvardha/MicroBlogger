-- +goose Up
create table books(
    id uuid not null primary key,
    name text not null,
    cover_image_url text not null,
    review text not null,
    tags text not null,
    level uuid not null references book_level(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table books;