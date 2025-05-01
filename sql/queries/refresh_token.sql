-- name: CreateRefreshToken :exec
insert into refresh_token(token, user_id, expires_at, created_at, updated_at)
values(
    $1,
    $2,
    $3,
    NOW(),
    NOW()
);

-- name: GetRefreshTokenExpirationTime :one
select expires_at from refresh_token where user_id = $1;

-- name: RemoveRefreshToken :exec
delete from refresh_token where token = $1 and user_id = $2;