create type FruitEnum as enum('apple', 'banana');

create table test_2_enum(
  fruit FruitEnum
)

