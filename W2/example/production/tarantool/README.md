
# see: https://kokizzu.blogspot.com/2021/05/easy-tarantool-clickhouse-replication-setup.html

## create table and insert on master
```shell
tarantoolctl connect tester:tester@127.0.0.1:13301
connected to 127.0.0.1:13301
```

```lua
127.0.0.1:13301> box.execute [[ create table test1(id int primary key, name string) ]]
---
- row_count: 1
...

127.0.0.1:13301> box.execute [[ insert into test1(id,name) values(1,'test') ]]
---
- row_count: 1
...
```

# check on slave cluster node
```shell
tarantoolctl connect tester:tester@127.0.0.1:23301
connected to 127.0.0.1:23301
```

```lua
127.0.0.1:23301> box.execute [[ select * FROM test1 ]]
---
- metadata:
  - name: ID
    type: integer
  - name: NAME
    type: string
  rows:
  - [1, 'test']
...
```
