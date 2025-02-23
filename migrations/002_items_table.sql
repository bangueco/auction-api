-- Write your migrate up statements here
create table items(
  id serial primary key,
  item_name varchar(255) not null,
  bid_amount float not null,
  auctioned_by serial references users(id)
);

---- create above / drop below ----
drop table items;
