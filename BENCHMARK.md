
# W2 Benchmark 2022-02-09

## Recap

| Name   | 10    | 1K    | 10K   | 20K   | p99 10 | p99 1K | p99  10K | p99  20K |
|--------|------:|------:|------:|------:|-------:|-------:|---------:|---------:|
| Write  | 24084 | 11505 | 11397 | 11112 | 0.0024 |   0.16 |     1.26 |     2.86 |
| Read   | 28596 | 10657 | 12296 | 10950 | 0.0007 |   0.17 |     1.30 |     2.78 |
| Health | 23568 | 11423 | 11195 | 11141 | 0.0008 |   0.16 |     1.71 |     3.37 |
| Hello  | 22780 | 12091 | 10365 | 10319 | 0.0009 |   0.16 |     1.84 |     3.98 |

Note:
- all benchmark done using [hey](//github.com/rakyll/hey) running 100K http requests, but with different concurrency levels on localhost
- logs enabled but discarded `make apiserver 2>&1 > /dev/null`, all dependencies run under `docker-compose`
- **Write** = insert into 1 table, retrieve the id, append the id to global array (ignoring race condition)
- **Read** = read random id from global array, then query from 1 table
- **Health** = do syscall (or read from /proc) cached once per second
- **Hello** = only serializing empty input and rendering `{"hello":"world"}` output
- Benchmark server: 32-core, 128GB RAM, NVMe disk

## C10 Write (database write, network call)

```shell
 hey -n 100000 -c 10 http://localhost:9090/api/LoadTestWrite

Summary:
  Total:        4.1521 secs
  Slowest:      0.0516 secs
  Fastest:      0.0001 secs
  Average:      0.0004 secs
  Requests/sec: 24084.1362
  
  Total data:   7499984 bytes
  Size/request: 74 bytes

Response time histogram:
  0.000 [1]     |
  0.005 [99494] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.010 [354]   |
  0.016 [96]    |
  0.021 [18]    |
  0.026 [12]    |
  0.031 [12]    |
  0.036 [5]     |
  0.041 [6]     |
  0.046 [1]     |
  0.052 [1]     |


Latency distribution:
  10% in 0.0002 secs
  25% in 0.0002 secs
  50% in 0.0003 secs
  75% in 0.0004 secs
  90% in 0.0006 secs
  95% in 0.0007 secs
  99% in 0.0024 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.0516 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0004 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0020 secs
  resp wait:    0.0004 secs, 0.0001 secs, 0.0515 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0020 secs

Status code distribution:
  [200] 100000 responses
```

## C1K Write

```shell
hey -n 100000 -c 1000 http://localhost:9090/api/LoadTestWrite

Summary:
  Total:        8.6918 secs
  Slowest:      0.3116 secs
  Fastest:      0.0002 secs
  Average:      0.0853 secs
  Requests/sec: 11505.1429
  
  Total data:   7500000 bytes
  Size/request: 75 bytes

Response time histogram:
  0.000 [1]     |
  0.031 [803]   |■
  0.063 [16525] |■■■■■■■■■■■■■
  0.094 [49475] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.125 [26771] |■■■■■■■■■■■■■■■■■■■■■■
  0.156 [5218]  |■■■■
  0.187 [733]   |■
  0.218 [380]   |
  0.249 [71]    |
  0.280 [17]    |
  0.312 [6]     |


Latency distribution:
  10% in 0.0557 secs
  25% in 0.0680 secs
  50% in 0.0830 secs
  75% in 0.0998 secs
  90% in 0.1171 secs
  95% in 0.1290 secs
  99% in 0.1624 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0006 secs, 0.0002 secs, 0.3116 secs
  DNS-lookup:   0.0009 secs, 0.0000 secs, 0.1909 secs
  req write:    0.0001 secs, 0.0000 secs, 0.0971 secs
  resp wait:    0.0839 secs, 0.0002 secs, 0.2249 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0809 secs

Status code distribution:
  [200] 100000 responses
```

## C10K Write

```shell
hey -n 100000 -c 10000 http://localhost:9090/api/LoadTestWrite

Summary:
  Total:        8.7736 secs
  Slowest:      1.7372 secs
  Fastest:      0.0064 secs
  Average:      0.8234 secs
  Requests/sec: 11397.8483
  
  Total data:   7500000 bytes
  Size/request: 75 bytes

Response time histogram:
  0.006 [1]     |
  0.179 [255]   |
  0.353 [670]   |■
  0.526 [3658]  |■■■■
  0.699 [17333] |■■■■■■■■■■■■■■■■■
  0.872 [41010] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.045 [26587] |■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.218 [8682]  |■■■■■■■■
  1.391 [1552]  |■■
  1.564 [227]   |
  1.737 [25]    |


Latency distribution:
  10% in 0.6174 secs
  25% in 0.7141 secs
  50% in 0.8183 secs
  75% in 0.9325 secs
  90% in 1.0507 secs
  95% in 1.1263 secs
  99% in 1.2666 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0233 secs, 0.0064 secs, 1.7372 secs
  DNS-lookup:   0.0162 secs, 0.0000 secs, 0.3918 secs
  req write:    0.0017 secs, 0.0000 secs, 0.3625 secs
  resp wait:    0.7867 secs, 0.0003 secs, 1.7371 secs
  resp read:    0.0004 secs, 0.0000 secs, 0.2636 secs

Status code distribution:
  [200] 100000 responses
```

## C20K Write

```shell
hey -n 100000 -c 20000 http://localhost:9090/api/LoadTestWrite

Summary:
  Total:        8.9988 secs
  Slowest:      3.8714 secs
  Fastest:      0.0008 secs
  Average:      1.6022 secs
  Requests/sec: 11112.6331
  
  Total data:   7500000 bytes
  Size/request: 75 bytes

Response time histogram:
  0.001 [1]     |
  0.388 [969]   |■
  0.775 [3177]  |■■■■
  1.162 [11032] |■■■■■■■■■■■■■
  1.549 [32934] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.936 [28344] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  2.323 [16345] |■■■■■■■■■■■■■■■■■■■■
  2.710 [5826]  |■■■■■■■
  3.097 [861]   |■
  3.484 [453]   |■
  3.871 [58]    |


Latency distribution:
  10% in 1.0419 secs
  25% in 1.3020 secs
  50% in 1.5702 secs
  75% in 1.9105 secs
  90% in 2.2135 secs
  95% in 2.4044 secs
  99% in 2.8638 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0619 secs, 0.0008 secs, 3.8714 secs
  DNS-lookup:   0.0965 secs, 0.0000 secs, 1.4742 secs
  req write:    0.0122 secs, 0.0000 secs, 0.5623 secs
  resp wait:    1.3670 secs, 0.0003 secs, 3.8000 secs
  resp read:    0.0168 secs, 0.0000 secs, 0.6904 secs

Status code distribution:
  [200] 100000 responses
```

## C 10 Read (database query, network call)

```shell
hey -n 100000 -c 10 http://localhost:9090/api/LoadTestRead 

Summary:
  Total:        3.4969 secs
  Slowest:      0.0027 secs
  Fastest:      0.0001 secs
  Average:      0.0003 secs
  Requests/sec: 28596.6464
  
  Total data:   37400000 bytes
  Size/request: 374 bytes

Response time histogram:
  0.000 [1]     |
  0.000 [58060] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.001 [36139] |■■■■■■■■■■■■■■■■■■■■■■■■■
  0.001 [5644]  |■■■■
  0.001 [119]   |
  0.001 [17]    |
  0.002 [8]     |
  0.002 [5]     |
  0.002 [0]     |
  0.002 [0]     |
  0.003 [7]     |


Latency distribution:
  10% in 0.0002 secs
  25% in 0.0002 secs
  50% in 0.0003 secs
  75% in 0.0005 secs
  90% in 0.0006 secs
  95% in 0.0006 secs
  99% in 0.0007 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.0027 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0004 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0005 secs
  resp wait:    0.0003 secs, 0.0001 secs, 0.0027 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0005 secs

Status code distribution:
  [200] 100000 responses
```

## C1K Read

```shell
hey -n 100000 -c 1000 http://localhost:9090/api/LoadTestRead 

Summary:
  Total:        9.3835 secs
  Slowest:      0.2796 secs
  Fastest:      0.0002 secs
  Average:      0.0921 secs
  Requests/sec: 10657.0083
  
  Total data:   37400000 bytes
  Size/request: 374 bytes

Response time histogram:
  0.000 [1]     |
  0.028 [544]   |■
  0.056 [8063]  |■■■■■■■■■
  0.084 [34160] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.112 [35047] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.140 [15844] |■■■■■■■■■■■■■■■■■■
  0.168 [4609]  |■■■■■
  0.196 [1415]  |■■
  0.224 [268]   |
  0.252 [35]    |
  0.280 [14]    |


Latency distribution:
  10% in 0.0579 secs
  25% in 0.0717 secs
  50% in 0.0890 secs
  75% in 0.1090 secs
  90% in 0.1302 secs
  95% in 0.1449 secs
  99% in 0.1776 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0003 secs, 0.0002 secs, 0.2796 secs
  DNS-lookup:   0.0010 secs, 0.0000 secs, 0.1742 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0754 secs
  resp wait:    0.0908 secs, 0.0002 secs, 0.2796 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0620 secs

Status code distribution:
  [200] 100000 responses
```

## C10K Read

```shell
hey -n 100000 -c 10000 http://localhost:9090/api/LoadTestRead 

Summary:
  Total:        8.1324 secs
  Slowest:      1.9896 secs
  Fastest:      0.0002 secs
  Average:      0.7559 secs
  Requests/sec: 12296.5247
  
  Total data:   37400000 bytes
  Size/request: 374 bytes

Response time histogram:
  0.000 [1]     |
  0.199 [841]   |■
  0.398 [3544]  |■■■
  0.597 [15991] |■■■■■■■■■■■■■■■■
  0.796 [40506] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.995 [26058] |■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.194 [10202] |■■■■■■■■■■
  1.393 [2522]  |■■
  1.592 [299]   |
  1.791 [32]    |
  1.990 [4]     |


Latency distribution:
  10% in 0.5123 secs
  25% in 0.6232 secs
  50% in 0.7409 secs
  75% in 0.8836 secs
  90% in 1.0360 secs
  95% in 1.1305 secs
  99% in 1.3002 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0213 secs, 0.0002 secs, 1.9896 secs
  DNS-lookup:   0.0089 secs, 0.0000 secs, 0.3246 secs
  req write:    0.0022 secs, 0.0000 secs, 0.3965 secs
  resp wait:    0.7148 secs, 0.0002 secs, 1.9896 secs
  resp read:    0.0019 secs, 0.0000 secs, 0.3288 secs
```

## C20K Read

```shell
hey -n 100000 -c 20000 http://localhost:9090/api/LoadTestRead 

Summary:
  Total:        9.1321 secs
  Slowest:      4.1688 secs
  Fastest:      0.0349 secs
  Average:      1.5950 secs
  Requests/sec: 10950.4404
  
  Total data:   37400000 bytes
  Size/request: 374 bytes

Response time histogram:
  0.035 [1]     |
  0.448 [1175]  |■
  0.862 [7044]  |■■■■■■■■
  1.275 [20223] |■■■■■■■■■■■■■■■■■■■■■■■■
  1.688 [33160] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  2.102 [19626] |■■■■■■■■■■■■■■■■■■■■■■■■
  2.515 [11557] |■■■■■■■■■■■■■■
  2.929 [6720]  |■■■■■■■■
  3.342 [346]   |
  3.755 [128]   |
  4.169 [20]    |


Latency distribution:
  10% in 0.9166 secs
  25% in 1.2202 secs
  50% in 1.5413 secs
  75% in 1.9559 secs
  90% in 2.4363 secs
  95% in 2.6717 secs
  99% in 2.7868 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0351 secs, 0.0349 secs, 4.1688 secs
  DNS-lookup:   0.0863 secs, 0.0000 secs, 1.5161 secs
  req write:    0.0700 secs, 0.0000 secs, 1.0982 secs
  resp wait:    1.3713 secs, 0.0002 secs, 3.6377 secs
  resp read:    0.0302 secs, 0.0000 secs, 1.2161 secs

Status code distribution:
  [200] 100000 responses
```

## C 10 Health (no database, syscall once per second)

```shell
 hey -n 100000 -c 10 http://localhost:9090/api/Health      

Summary:
  Total:        4.2430 secs
  Slowest:      0.0044 secs
  Fastest:      0.0000 secs
  Average:      0.0004 secs
  Requests/sec: 23568.3716
  
  Total data:   9699940 bytes
  Size/request: 96 bytes

Response time histogram:
  0.000 [1]     |
  0.000 [62082] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.001 [37697] |■■■■■■■■■■■■■■■■■■■■■■■■
  0.001 [189]   |
  0.002 [7]     |
  0.002 [4]     |
  0.003 [0]     |
  0.003 [0]     |
  0.004 [4]     |
  0.004 [1]     |
  0.004 [15]    |


Latency distribution:
  10% in 0.0002 secs
  25% in 0.0003 secs
  50% in 0.0004 secs
  75% in 0.0005 secs
  90% in 0.0006 secs
  95% in 0.0007 secs
  99% in 0.0008 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0044 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0003 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0011 secs
  resp wait:    0.0004 secs, 0.0000 secs, 0.0044 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0012 secs

Status code distribution:
  [200] 100000 responses
```

## C1K Health

```shell
hey -n 100000 -c 1000 http://localhost:9090/api/Health 

Summary:
  Total:        8.7541 secs
  Slowest:      0.3141 secs
  Fastest:      0.0001 secs
  Average:      0.0858 secs
  Requests/sec: 11423.2743
  
  Total data:   9887748 bytes
  Size/request: 98 bytes

Response time histogram:
  0.000 [1]     |
  0.031 [914]   |■
  0.063 [17483] |■■■■■■■■■■■■■■■
  0.094 [47910] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.126 [26679] |■■■■■■■■■■■■■■■■■■■■■■
  0.157 [5814]  |■■■■■
  0.188 [763]   |■
  0.220 [396]   |
  0.251 [18]    |
  0.283 [14]    |
  0.314 [8]     |


Latency distribution:
  10% in 0.0550 secs
  25% in 0.0678 secs
  50% in 0.0835 secs
  75% in 0.1012 secs
  90% in 0.1194 secs
  95% in 0.1311 secs
  99% in 0.1612 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0004 secs, 0.0001 secs, 0.3141 secs
  DNS-lookup:   0.0009 secs, 0.0000 secs, 0.1927 secs
  req write:    0.0001 secs, 0.0000 secs, 0.1903 secs
  resp wait:    0.0846 secs, 0.0001 secs, 0.2442 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0686 secs

Status code distribution:
  [200] 100000 responses
```

## C10K Health

```shell
hey -n 100000 -c 10000 http://localhost:9090/api/Health

Summary:
  Total:        8.9323 secs
  Slowest:      3.1181 secs
  Fastest:      0.0001 secs
  Average:      0.8167 secs
  Requests/sec: 11195.3508
  
  Total data:   9854998 bytes
  Size/request: 98 bytes

Response time histogram:
  0.000 [1]     |
  0.312 [2875]  |■■■
  0.624 [24787] |■■■■■■■■■■■■■■■■■■■■■■■■
  0.936 [41571] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.247 [21197] |■■■■■■■■■■■■■■■■■■■■
  1.559 [7329]  |■■■■■■■
  1.871 [1812]  |■■
  2.183 [362]   |
  2.495 [55]    |
  2.806 [10]    |
  3.118 [1]     |


Latency distribution:
  10% in 0.4914 secs
  25% in 0.6094 secs
  50% in 0.7559 secs
  75% in 0.9941 secs
  90% in 1.2371 secs
  95% in 1.3970 secs
  99% in 1.7137 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0137 secs, 0.0001 secs, 3.1181 secs
  DNS-lookup:   0.0084 secs, 0.0000 secs, 0.3815 secs
  req write:    0.0014 secs, 0.0000 secs, 0.3469 secs
  resp wait:    0.7843 secs, 0.0001 secs, 3.1180 secs
  resp read:    0.0010 secs, 0.0000 secs, 0.3027 secs

Status code distribution:
  [200] 100000 responses
```

## C20K Health

```shell
hey -n 100000 -c 20000 http://localhost:9090/api/Health      

Summary:
  Total:        8.9757 secs
  Slowest:      5.5716 secs
  Fastest:      0.0001 secs
  Average:      1.5825 secs
  Requests/sec: 11141.2181
  
  Total data:   9877214 bytes
  Size/request: 98 bytes

Response time histogram:
  0.000 [1]     |
  0.557 [4135]  |■■■■
  1.114 [13631] |■■■■■■■■■■■■
  1.672 [44756] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  2.229 [21275] |■■■■■■■■■■■■■■■■■■■
  2.786 [11910] |■■■■■■■■■■■
  3.343 [3163]  |■■■
  3.900 [1003]  |■
  4.457 [111]   |
  5.014 [13]    |
  5.572 [2]     |


Latency distribution:
  10% in 0.9157 secs
  25% in 1.2031 secs
  50% in 1.4036 secs
  75% in 1.9678 secs
  90% in 2.4978 secs
  95% in 2.7296 secs
  99% in 3.3747 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0684 secs, 0.0001 secs, 5.5716 secs
  DNS-lookup:   0.1299 secs, 0.0000 secs, 2.4811 secs
  req write:    0.0115 secs, 0.0000 secs, 1.8493 secs
  resp wait:    1.3601 secs, 0.0001 secs, 4.5270 secs
  resp read:    0.0213 secs, 0.0000 secs, 1.8238 secs

Status code distribution:
  [200] 100000 responses
```

## C 10 Hello (serialization-deserialization only)

```shell
hey -n 100000 -c 10 http://localhost:9090/api/LoadHello

Summary:
  Total:        4.3897 secs
  Slowest:      0.0226 secs
  Fastest:      0.0000 secs
  Average:      0.0004 secs
  Requests/sec: 22780.7162
  
  Total data:   5700000 bytes
  Size/request: 57 bytes

Response time histogram:
  0.000 [1]     |
  0.002 [99959] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.005 [29]    |
  0.007 [0]     |
  0.009 [0]     |
  0.011 [1]     |
  0.014 [0]     |
  0.016 [0]     |
  0.018 [0]     |
  0.020 [0]     |
  0.023 [10]    |


Latency distribution:
  10% in 0.0002 secs
  25% in 0.0003 secs
  50% in 0.0004 secs
  75% in 0.0006 secs
  90% in 0.0007 secs
  95% in 0.0007 secs
  99% in 0.0009 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0226 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0003 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0036 secs
  resp wait:    0.0004 secs, 0.0000 secs, 0.0225 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0077 secs

Status code distribution:
  [200] 100000 responses
```

## C1K Hello

```shell
hey -n 100000 -c 1000 http://localhost:9090/api/LoadHello 

Summary:
  Total:        8.2700 secs
  Slowest:      0.3034 secs
  Fastest:      0.0001 secs
  Average:      0.0813 secs
  Requests/sec: 12091.9439
  
  Total data:   5700000 bytes
  Size/request: 57 bytes

Response time histogram:
  0.000 [1]     |
  0.030 [1078]  |■
  0.061 [19059] |■■■■■■■■■■■■■■■
  0.091 [49452] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.121 [24805] |■■■■■■■■■■■■■■■■■■■■
  0.152 [4347]  |■■■■
  0.182 [367]   |
  0.212 [479]   |
  0.243 [299]   |
  0.273 [96]    |
  0.303 [17]    |


Latency distribution:
  10% in 0.0518 secs
  25% in 0.0641 secs
  50% in 0.0788 secs
  75% in 0.0953 secs
  90% in 0.1120 secs
  95% in 0.1233 secs
  99% in 0.1641 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0006 secs, 0.0001 secs, 0.3034 secs
  DNS-lookup:   0.0010 secs, 0.0000 secs, 0.1533 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0798 secs
  resp wait:    0.0798 secs, 0.0001 secs, 0.2400 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0450 secs

Status code distribution:
  [200] 100000 responses
```

## C10K Hello

```shell
hey -n 100000 -c 10000 http://localhost:9090/api/LoadHello

Summary:
  Total:        9.6475 secs
  Slowest:      2.9423 secs
  Fastest:      0.0001 secs
  Average:      0.8973 secs
  Requests/sec: 10365.3998
  
  Total data:   5700000 bytes
  Size/request: 57 bytes

Response time histogram:
  0.000 [1]     |
  0.294 [2009]  |■■
  0.589 [11498] |■■■■■■■■■■■
  0.883 [40739] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.177 [27869] |■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.471 [12220] |■■■■■■■■■■■■
  1.765 [4213]  |■■■■
  2.060 [1138]  |■
  2.354 [242]   |
  2.648 [59]    |
  2.942 [12]    |


Latency distribution:
  10% in 0.5442 secs
  25% in 0.6792 secs
  50% in 0.8466 secs
  75% in 1.0933 secs
  90% in 1.3220 secs
  95% in 1.4980 secs
  99% in 1.8467 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0165 secs, 0.0001 secs, 2.9423 secs
  DNS-lookup:   0.0142 secs, 0.0000 secs, 0.2916 secs
  req write:    0.0009 secs, 0.0000 secs, 0.1381 secs
  resp wait:    0.8636 secs, 0.0001 secs, 2.9423 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0806 secs

Status code distribution:
  [200] 100000 responses
```

## C20K Hello

```shell
hey -n 100000 -c 20000 http://localhost:9090/api/LoadHello

Summary:
  Total:        9.6904 secs
  Slowest:      6.4666 secs
  Fastest:      0.0002 secs
  Average:      1.6755 secs
  Requests/sec: 10319.4446
  
  Total data:   5700000 bytes
  Size/request: 57 bytes

Response time histogram:
  0.000 [1]     |
  0.647 [5335]  |■■■■■■
  1.293 [33214] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.940 [27625] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  2.587 [21133] |■■■■■■■■■■■■■■■■■■■■■■■■■
  3.233 [8663]  |■■■■■■■■■■
  3.880 [2849]  |■■■
  4.527 [869]   |■
  5.173 [261]   |
  5.820 [39]    |
  6.467 [11]    |


Latency distribution:
  10% in 0.8219 secs
  25% in 1.1397 secs
  50% in 1.5045 secs
  75% in 2.1292 secs
  90% in 2.7457 secs
  95% in 3.1096 secs
  99% in 3.9815 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0575 secs, 0.0002 secs, 6.4666 secs
  DNS-lookup:   0.0340 secs, 0.0000 secs, 0.6683 secs
  req write:    0.0233 secs, 0.0000 secs, 0.8265 secs
  resp wait:    1.5186 secs, 0.0001 secs, 6.4665 secs
  resp read:    0.0029 secs, 0.0000 secs, 0.8066 secs

Status code distribution:
  [200] 100000 responses
```
