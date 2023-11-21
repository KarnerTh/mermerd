-- Test case for https://github.com/KarnerTh/mermerd/issues/15
-- only stub, because sqlite does not support enums
create table test_2_enum(
  fruit TEXT CHECK( fruit IN ('apple','banana') )
);

