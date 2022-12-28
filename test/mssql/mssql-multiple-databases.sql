-- Test case for https://github.com/KarnerTh/mermerd/issues/23
create table dbo.test_3_a
(
    id    int          not null primary key,
    title varchar(255) not null
);

GO

create schema other_db;

GO

create table other_db.test_3_b
(
    id    int          not null primary key,
    aid int,
    foreign key (aid) references dbo.test_3_a (id)
);

create table other_db.test_3_c
(
    id    int          not null primary key,
    title varchar(255) not null
);

