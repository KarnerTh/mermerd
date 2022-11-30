create table article
(
    id    int          not null primary key,
    title varchar(255) not null
);

create table article_detail
(
    id         int       not null primary key,
    created_at timestamp not null
);

create table article_comment
(
    id         int          not null primary key,
    article_id int          not null,
    comment    varchar(255) not null
);

create table label
(
    id    int          not null primary key,
    label varchar(255) not null
);

create table article_label
(
    article_id int not null,
    label_id   int not null
);

-- one-to-one relation
alter table article_detail
    add constraint fk_article_detail_id foreign key (id) references article (id);

-- one-to-many relation
alter table article_comment
    add constraint fk_article_comment_id foreign key (article_id) references article (id);

-- many-to-many relation
alter table article_label
    add constraint fk_article_label_article_id foreign key (article_id) references article (id);
alter table article_label
    add constraint fk_article_label_label_id foreign key (label_id) references label (id);
alter table article_label
    add primary key (article_id, label_id);

-- Test case for https://github.com/KarnerTh/mermerd/issues/8
create table test_1_a
(
    id  int,
    xid int,
    primary key (id, xid)
);

create table test_1_b
(
    aid int,
    bid int,
    primary key (aid, bid),
    foreign key (aid, bid) references test_1_a (id, xid)
);

