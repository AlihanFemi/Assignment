-- +goose Up
-- +goose StatementBegin
create table comments(
	ID serial primary key,
	content varchar(512) not null,
	Date varchar(10) not null,
	post_ID int not null,
	user_ID int not null
);
create table posts(
	ID serial primary key,
	Title varchar(64) not null,
	content varchar(1024) not null,
	Date varchar(10) not null,
	user_ID int not null
);
create table users(
	ID serial primary key,
	Name varchar(64) not null,
	Email varchar(64) not null,
	Password varchar(72) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
