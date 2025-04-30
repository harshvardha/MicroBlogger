-- name: CreateBlog :one
insert into blogs(
    id, title, brief, content_url,
    images, thumbnail_url, code_repo_link, tags,
    author, category, created_at, updated_at
) values(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateBlog :exec
update blogs set
title = $1, brief = $2, content_url = $3,
thumbnail_url = $4, code_repo_link = $5,
category = $6, updated_at = NOW() where id = $7;

-- name: UpdateImages :exec
update blogs set images = $1, updated_at = NOW() where id = $2;

-- name: UpdateBlogTags :exec
update blogs set tags = $1, updated_at = NOW() where id = $2;

-- name: RemoveBlog :exec
delete from blogs where id = $1;

-- name: GetBlogByID :one
select
blogs.title, blogs.content_url, blogs.images,
blogs.code_repo_link, blogs.views, blogs.likes,
blogs.tags, users.username, blogs.created_at
from blogs join users on blogs.author = users.id where blogs.id = $1;

-- name: GetAllBlogsByCategory :many
select 
id, title, brief,
thumbnail_url, views, likes,
tags, created_at from blogs where category = $1 and id > $2 limit $3;