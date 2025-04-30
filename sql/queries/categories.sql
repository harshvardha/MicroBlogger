-- name: CreateCategory :one
insert into categories(id, category, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateCategory :one
update categories set category = $1, updated_at = NOW() where id = $2
returning category, updated_at;

-- name: RemoveCategory :exec
delete from categories where id = $1;

-- name: GetAllCategories :many
select * from categories;