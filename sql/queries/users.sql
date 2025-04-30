-- name: CreateUser :one
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
)
returning id, username, profile_pic_url, role_id, created_at, updated_at;

-- name: UpdateUsername :one
update users set username = $1, updated_at = NOW() where id = $2
returning username;

-- name: UpdateProfilePic :one
update users set profile_pic_url = $1, updated_at = NOW() where id = $2
returning profile_pic_url;

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