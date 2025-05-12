-- +goose Up
create table likes(
    user_id uuid not null references users(id) on delete cascade,
    blog_id uuid not null references blogs(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null,
    unique(user_id, blog_id)
);

-- +goose Down
drop table likes;