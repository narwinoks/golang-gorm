CREATE TABLE users (
    id varchar(100) not null,
    name varchar(100) not null,
    password varchar(255) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,

    primary key (id)
) engine  =InnoDB;


# DROP table users

alter table users rename column name to first_name;

alter table users add column middle_name  varchar(100) null after first_name;
alter table users add column last_name varchar(100) null after middle_name;

select * from users;
delete from users where id = 2;


create table user_logs(
    id  int auto_increment,
    user_id varchar(100) not null,
    action varchar(100) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key (id)
) engine  =InnoDB;


select * from user_logs;

delete from user_logs;

ALTER table user_logs modify created_at bigint not null;
ALTER table user_logs modify updated_at bigint not null;


select * from user_logs;

create table todos (
    id bigint not null  auto_increment,
    user_id varchar( 100 ) not null,
    title varchar(100) not null,
    description text null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    deleted_at timestamp null,

    primary key  (id)
)engine  = innodb;


ALTER TABLE todos RENAME COLUMN created_t TO  created_at;

    select * from todos;

create table wallets(
    id varchar(100) not null,
    user_id varchar(100) not null,
    balance bigint not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp on update current_timestamp,
    primary key  (id),
    unique (user_id)
)engine = innodb;


desc wallets;

create table addresses(
    id bigint not null auto_increment,
    user_id varchar(100) not null,
    address varchar(100) not null,
    created_at timestamp not null default current_timestamp,
    update_at timestamp not null default  current_timestamp on update current_timestamp,
    primary key(id),
    foreign key (user_id) references users(id)
)engine  = innodb;
ALTER TABLE addresses RENAME COLUMN update_at TO  updated_at;
desc
