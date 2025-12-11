-- Minimal test case for cross-schema FK bug
-- Bug: GetConstraints() fails when multiple schemas have identically-named constraints
-- Reproducer: Two schemas both have "users_pkey", cross-schema FK between them

create schema tenant_1;
create table tenant_1.users
(
    id int constraint users_pkey primary key
);

create schema tenant_2;
create table tenant_2.users
(
    id int constraint users_pkey primary key
);


create table tenant_2.posts
(
    id         int primary key,
    author_id  int references tenant_1.users(id)
);
