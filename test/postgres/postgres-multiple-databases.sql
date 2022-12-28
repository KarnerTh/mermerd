-- Test case for https://github.com/KarnerTh/mermerd/issues/23
create table test_3_a
(
    id    int          not null primary key,
    title varchar(255) not null
);

create schema other_db;

create table other_db.test_3_b
(
    id    int          not null primary key,
    aid   int,
    foreign key (aid) references public.test_3_a (id)
);

create table other_db.test_3_c
(
    id    int          not null primary key,
    title varchar(255) not null
);

