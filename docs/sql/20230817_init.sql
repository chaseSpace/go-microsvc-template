create database biz;
create database biz_log;

drop table if exists biz.user;
create table biz.user
(
    id          int primary key auto_increment,
    uid         int              not null,
    ext_uid     int              not null,
    nickname    varchar(20)      not null,
    age         tinyint unsigned not null,
    sex         tinyint unsigned not null,
    account     varchar(30)      not null,
    password    varchar(40)      not null comment 'encrypted with salt in SHA1',
    passwd_salt varchar(10)      not null,
    created_at  datetime         not null default current_timestamp,
    updated_at  datetime         not null default current_timestamp on update current_timestamp,
    key idx_uid (uid),
    key idx_extuid (ext_uid),
    key idx_ct (created_at)
)
