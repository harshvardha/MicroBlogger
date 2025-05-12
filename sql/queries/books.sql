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

-- name: UpdateBook :exec
update books
set name = $1, cover_image_url = $2,
review = $3, tags = $4, level = $5, updated_at = NOW()
where id = $6;

-- name: RemoveBook :exec
delete from books where id = $1;

-- name: GetBooksByLevel :many
select id, name, cover_image_url from books where books.level = (select id from book_level where book_level.level = $1);

-- name: GetAllBooks :many
select id, name, cover_image_url from books;

-- name: GetAllBooksCount :one
select count(*) from books;

-- name: GetLevelIDByName :one
select id from book_level where level = $1;

-- name: GetReviewByBookID :one
select review, cover_image_url from books where id = $1;

-- name: GetBookByID :one
select name, cover_image_url, review, tags, level from books where id = $1;