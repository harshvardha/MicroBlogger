-- name: CreateRole :one
insert into roles(id, role_name, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    NOW(),
    NOW()
)
returning *;

-- name: RemoveRole :exec
delete from roles where id = $1;

-- name: GetRoleIdByName :one
select id from roles where role_name = $1;

-- name: GetRoleById :one
select * from roles where id = $1;

-- name: UpdateRoleById :one
update roles set role_name = $1, updated_at = NOW() where id = $2
returning role_name;