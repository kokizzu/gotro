

# see: https://kokizzu.blogspot.com/2021/05/easy-tarantool-clickhouse-replication-setup.html

## create table and insert on first cluster node

```shell
clickhouse-client --port 19000
```

```sql
SELECT * FROM system.clusters;
CREATE DATABASE db1 ON CLUSTER replicated;
SHOW DATABASES;
USE db1;

CREATE TABLE IF NOT EXISTS db1.table1 ON CLUSTER replicated
( id UInt64
, dt Date
, val UInt64
) ENGINE = ReplicatedReplacingMergeTree('/clickhouse/{cluster}/tables/table1',
'{replica}')
PARTITION BY modulo( id, 1000 )
ORDER BY (dt);

INSERT INTO db1.table1
(id, dt, val)
VALUES (1,'2021-05-31',2);
```

## check on second cluster node
```shell
clickhouse-client --port 29000
```

```sql
SELECT * FROM db1.sbr2;

┌─id─┬─dt─────────┬─val─┐
│  1 │ 2021-05-31 │   2 │
└────┴────────────┴─────┘
↘ Progress: 1.00 rows, 42.00 B (132.02 rows/s., 5.54 KB/s.)  99%
1 rows in set. Elapsed: 0.008 sec.
```
