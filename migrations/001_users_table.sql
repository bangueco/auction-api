-- Write your migrate up statements here
create table users(
  id serial primary key,
  username varchar(255) not null unique,
  password varchar(255) not null
);

---- create above / drop below ----

drop table users;
