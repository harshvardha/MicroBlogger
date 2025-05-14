-- +goose Up
alter table blogs drop column likes;

-- +goose Down
alter table blogs add column likes int;