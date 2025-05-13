-- +goose Up
alter table blogs add column tags_new text[] not null;
alter table books add column tags_new text[] not null;
alter table blogs drop column tags;
alter table books drop column tags;
alter table blogs rename column tags_new to tags;
alter table books rename column tags_new to tags;

-- +goose Down
alter table blogs add column tags_old text not null;
alter table books add column tags_old text not null;
update table blogs set tags_old = array_to_string(tags, ';');
update table books set tags_old = array_to_string(tags, ';');
alter table blogs drop column tags;
alter table books drop column tags;
alter table blogs rename column tags_old to tags;
alter table books rename column tags_old to tags;