#!/bin/bash
# Source: https://cardano.github.io/blog/2017/11/15/mssql-docker-container
wait_time=15s
password=securePassword1!

# wait for SQL Server to come up
echo importing data will start in $wait_time...
sleep $wait_time
echo importing data...

# run the sql scripts to create the test database
/opt/mssql-tools/bin/sqlcmd -S 0.0.0.0 -U sa -P $password -i ./mssql-setup.sql
/opt/mssql-tools/bin/sqlcmd -S 0.0.0.0 -U sa -P $password -d mermerd_test -i ./db-table-setup.sql
/opt/mssql-tools/bin/sqlcmd -S 0.0.0.0 -U sa -P $password -d mermerd_test -i ./mssql-enum-setup.sql
/opt/mssql-tools/bin/sqlcmd -S 0.0.0.0 -U sa -P $password -d mermerd_test -i ./mssql-multiple-databases.sql

echo importing done
