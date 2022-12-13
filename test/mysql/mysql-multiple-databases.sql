-- Test case for https://github.com/KarnerTh/mermerd/issues/23
create table test_3_a
(
    id    int          not null primary key,
    title varchar(255) not null
);

create database other_db;
use other_db;

create table test_3_b
(
    id    int          not null primary key,
    aid int,
    foreign key (aid) references mermerd_test.test_3_a (id)
);

create table test_3_c
(
    id    int          not null primary key,
    title varchar(255) not null
);

