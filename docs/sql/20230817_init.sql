create database microsvc;

create table microsvc.user
(
    id         int primary key auto_increment,
    uid        int              not null,
    alias_uid  int              not null comment '对外展示的UID',
    nick       varchar(50)      not null,
    age        tinyint unsigned not null,
    sex        tinyint unsigned not null,
    created_at datetime         not null default current_timestamp,
    updated_at datetime         not null default current_timestamp on update current_timestamp,
    key idx_uid (uid),
    key idx_aliasuid (alias_uid),
    key idx_ct (created_at)
)