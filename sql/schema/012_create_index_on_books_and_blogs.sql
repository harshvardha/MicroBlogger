-- +goose Up
create extension if not exists pg_trgm;
create index idx_blogs_title_trgm on blogs using gin(title gin_trgm_ops);
create index idx_blogs_title_prefix on blogs(title text_pattern_ops);
create index idx_blogs_tags on blogs using gin(tags);
create index idx_books_name_trgm on books using gin(name gin_trgm_ops);
create index idx_books_tags on books using gin(tags);

-- +goose Down
drop extension if exists pg_trgm;
drop index if exists idx_blogs_title_trgm;
drop index if exists idx_blogs_title_prefix;
drop index if exists idx_blogs_tags;
drop index if exists idx_books_name_trgm;
drop index if exists idx_books_tags;