-- name: CreateBook :one
insert into books(
    id, name, cover_image_url,review,
    tags, level, created_at, updated_at
) values(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateBook :one
update books
set name = $1, cover_image_url = $2,
review = $3, level = $4, updated_at = NOW()
where id = $5
returning *;

-- name: UpdateBookTags :one
update books set tags = $1 where id = $2
returning tags;

-- name: RemoveBook :exec
delete from books where id = $1;

-- name: GetBooksByLevel :many
select name, cover_image_url, review, tags from books where level = $1;

-- name: GetAllBooks :many
select * from books;

-- name: GetAllBooksCount :one
select count(*) from books;