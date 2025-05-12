-- +goose Up
create table blogs(
    id uuid not null primary key,
    title text not null unique,
    brief varchar(200) not null unique,
    content_url text not null,
    images json not null,
    thumbnail_url text not null,
    code_repo_link text,
    views int not null default 0,
    likes int not null default 0,
    tags text not null,
    author uuid not null references users(id) on delete cascade,
    category uuid not null references categories(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table blogs;