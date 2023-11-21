create table article
(
    id    int          not null primary key,
    title varchar(255) not null,
    subtitle varchar(255)
);

create table article_detail
(
    id         int       not null primary key,
    created_at timestamp not null,
    foreign key (id) references article (id)
);

create table article_comment
(
    id         int          not null primary key,
    article_id int          not null,
    comment    varchar(255) not null,
    foreign key (article_id) references article (id)
);

create table label
(
    id    int          not null primary key,
    label varchar(255) not null unique
);

create table article_label
(
    article_id int not null,
    label_id   int not null,
    foreign key (article_id) references article (id),
    foreign key (label_id) references label (id),
    primary key (article_id, label_id)
);

-- Test case for https://github.com/KarnerTh/mermerd/issues/8
create table test_1_a
(
    id  int not null,
    xid int not null,
    primary key (id, xid)
);

create table test_1_b
(
    aid int not null,
    bid int not null,
    primary key (aid, bid),
    foreign key (aid, bid) references test_1_a (id, xid)
);

