-- name: CreateComment :one
insert into comments(
    id, description, user_id, blog_id,
    created_at, updated_at
) values(
    gen_random_uuid(),
    $1, $2, $3, NOW(), NOW()
)
returning id, description, created_at, updated_at;

-- name: UpdateCommentByID :one
update comments set description = $1, updated_at = NOW() where id = $2 and user_id = $3
returning description, updated_at;

-- name: GetCommentByBlogID :many
select
comments.id, comments.description, users.username,
users.profile_pic_url, comments.created_at, comments.updated_at
from comments join users on comments.user_id = users.id where comments.blog_id = $1;

-- name: RemoveComment :exec
delete from comments where id = $1 and user_id = $2;