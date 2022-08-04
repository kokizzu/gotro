
-- clickhouse-client --port 19000
-- clickhouse-client --port 29000

CREATE DATABASE db1 ON CLUSTER replicated;

-- stores rocket-s3 put object event
CREATE TABLE IF NOT EXISTS db1.puts ON CLUSTER replicated
( root String
, sub String
, bucket String
, obj String
, dc String
, ts DateTime
, size UInt64
) ENGINE = ReplicatedReplacingMergeTree('/clickhouse/{cluster}/tables/puts',
'{replica}')
PARTITION BY modulo( xxHash32(root), 1000 )
ORDER BY (root,sub,bucket,obj,dc,ts);

-- stores rocket-s3 delete object event
CREATE TABLE IF NOT EXISTS db1.dels ON CLUSTER replicated
( root String
, sub String
, bucket String
, obj String
, dc String
, ts DateTime
, at DateTime -- GREATEST(at,addDays(ts,60)) -- 60 is min_object_billing_days
, size UInt64
) ENGINE = ReplicatedReplacingMergeTree('/clickhouse/{cluster}/tables/dels',
'{replica}')
PARTITION BY modulo( xxHash32(root), 1000 )
ORDER BY (root,sub,bucket,obj,dc,ts,at);

-- stores rwave successful copy event
CREATE TABLE IF NOT EXISTS db1.dups ON CLUSTER replicated
( root String
, sub String
, bucket String
, obj String
, src String -- datacenter
, dst String
, ts DateTime
, size UInt64
) ENGINE = ReplicatedReplacingMergeTree('/clickhouse/{cluster}/tables/dups',
'{replica}')
PARTITION BY modulo( xxHash32(root), 1000 )
ORDER BY (root,sub,bucket,obj,src,dst,ts);

-- stores rwave successful delete event
CREATE TABLE IF NOT EXISTS db1.wipes ON CLUSTER replicated
( root String
, sub String
, bucket String
, obj String
, dc String
, ts DateTime
, size UInt64
) ENGINE = ReplicatedReplacingMergeTree('/clickhouse/{cluster}/tables/wipes',
'{replica}')
PARTITION BY modulo( xxHash32(root), 1000 )
ORDER BY (root,sub,bucket,obj,dc,ts);

-- stores sum put per reseller
CREATE TABLE IF NOT EXISTS db1.sum_puts ON CLUSTER replicated
( root String
, size UInt64
) ENGINE = ReplicatedSummingMergeTree('/clickhouse/{cluster}/tables/sum_puts1',
'{replica}')
ORDER BY (root);

CREATE MATERIALIZED VIEW db1.mv_sum_puts ON CLUSTER replicated
TO db1.sum_puts AS 
  SELECT root, SUM(size) AS size
  FROM db1.puts
  GROUP BY root;

INSERT INTO db1.puts VALUES('reseller1','customer1','bucket1','object1','DC1',NOW(),1);
INSERT INTO db1.puts VALUES('reseller1','customer1','bucket1','object2','DC1',NOW(),2);

SELECT * FROM db1.sum_puts FINAL; -- 3

DROP VIEW db1.mv_sum_puts ON CLUSTER replicated SYNC;
DROP TABLE db1.sum_puts ON CLUSTER replicated SYNC;
