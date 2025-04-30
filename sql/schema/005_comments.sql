-- +goose Up
create table comments(
    id uuid not null primary key,
    description text not null,
    likes int not null default 0,
    user_id uuid not null references users(id) on delete cascade,
    blog_id uuid not null references blogs(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table comments;