create database biz;
create database biz_log;

drop table if exists biz.user;
create table biz.user
(
    id            int primary key auto_increment,
    uid           int              not null,
    nickname      varchar(20)      not null,
    birthday      date             not null,
    sex           tinyint unsigned not null comment '0-未知 1男2女',
    password      varchar(40)      not null,
    password_salt varchar(10)      not null,
    created_at    datetime         not null default current_timestamp,
    updated_at    datetime         not null default current_timestamp on update current_timestamp,
    unique key idx_uid (uid),
    key idx_ct (created_at)
)
