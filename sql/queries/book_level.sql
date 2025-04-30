-- name: CreateLevel :one
insert into book_level(id, level, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateLevel :one
update book_level set level = $1, updated_at = NOW() where id = $2
returning *;

-- name: RemoveLevel :exec
delete from book_level where id = $1;

-- name: GetAllBookLevels :many
select * from book_level;