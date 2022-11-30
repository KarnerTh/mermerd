-- Test case for https://github.com/KarnerTh/mermerd/issues/15
-- due to the fact that mssql does not support enums, we use the closest possible solution 
-- https://stackoverflow.com/a/1434338
create table test_2_enum(
  fruit varchar(10) not null check (fruit in('apple', 'banana'))
)

