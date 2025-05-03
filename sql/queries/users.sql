-- name: CreateUser :exec
insert into users(id, email, username, profile_pic_url, password, role_id, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    NOW(),
    NOW()
);

-- name: UpdateUsername :one
update users set username = $1, profile_pic_url = $2, updated_at = NOW() where id = $3
returning username;

-- name: UpdateEmail :exec
update users set email = $1, updated_at = NOW() where id = $2;

-- name: UpdatePassword :exec
update users set password = $1, updated_at = NOW() where id = $2;

-- name: RemoveUser :exec
delete from users where id = $1;

-- name: GetUserByID :one
select 
    users.email, 
    users.username, 
    users.profile_pic_url, 
    users.created_at, 
    users.updated_at,
    roles.role_name from users join roles on users.role_id = roles.id where users.id = $1;

-- name: GetUserRole :one
select roles.role_name from users join roles on users.id = roles.id where users.id = $1;

-- name: UserExist :one
select exists(select 1 from users where email = $1);

-- name: GetUserByEmailID :one
select id, username, profile_pic_url, password, role_id from users where email = $1;