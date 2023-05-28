-- Test case for https://github.com/KarnerTh/mermerd/issues/36
create table test_not_unique_constraint_name_a
(
    id int primary key
);
create table test_not_unique_constraint_name_b
(
    id int primary key
);
create table test_not_unique_constraint_name_c
(
    id int primary key
);

alter table test_not_unique_constraint_name_b 
    add constraint not_unique_constraint_name foreign key (id) references test_not_unique_constraint_name_a (id);
alter table test_not_unique_constraint_name_c 
    add constraint not_unique_constraint_name foreign key (id) references test_not_unique_constraint_name_a (id);
