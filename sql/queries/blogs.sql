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

-- name: UpdateBlog :one
update blogs set
title = $1, brief = $2, content_url = $3, images = $4,
thumbnail_url = $5, code_repo_link = $6, tags = $7,
category = $8, updated_at = NOW() where id = $9
returning created_at, updated_at;

-- name: RemoveBlog :exec
delete from blogs where id = $1;

-- name: GetBlogByID :one
select
blogs.title, blogs.brief, blogs.content_url, blogs.images, blogs.thumbnail_url,
blogs.code_repo_link, blogs.views, blogs.likes,
blogs.tags, users.username, blogs.created_at
from blogs join users on blogs.author = users.id where blogs.id = $1;

-- name: GetAllBlogsByCategory :many
select 
id, title, brief, thumbnail_url, views, likes,
tags, created_at from blogs where category = $1 and id > $2 limit $3;

-- name: LikeBlog :exec
insert into likes(user_id, blog_id, created_at, updated_at)
values($1, $2, NOW(), NOW());

-- name: DislikeBlog :exec
delete from likes where user_id = $1 and blog_id = $2;

-- name: GetNumberOfLikes :one
select count(*) as noOfLikes from likes where blog_id = $1;

-- name: HasUserLikedBlog :one
select 1 from likes where user_id = $1 and blog_id = $2;

-- name: IncrementViews :exec
update blogs set views = views + 1 where id = $1;

-- name: GetViewCount :one
select views from blogs where id = $1;