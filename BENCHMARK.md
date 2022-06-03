
# W2 Benchmark 2022-02-09

## Recap

| Name   |     10 |     1K |   10K |   20K | p99 10 | p99 1K | p99 10K | p99 20K |
|--------|-------:|------:|------:|------:|-------:|-------:|--------:|--------:|
| Write  |  46695 |  55719 | 52232 | 33428 | 0.0004 | 0.0656 |  0.3889 |  1.9080 |
| Read   |  54946 | 112834 | 82359 | 38736 | 0.0004 | 0.0761 |  0.5704 |  1.8054 |
| Health | 106305 | 143231 | 99251 | 51664 | 0.0003 | 0.0600 |  0.5148 |  1.4058 |
| Hello  | 104708 | 144292 | 98795 | 44347 | 0.0003 | 0.0389 |  0.4940 |  1.6238 |

Note:
- all benchmark done using [hey](//github.com/rakyll/hey) running 100K http requests, but with different concurrency levels on localhost
- logs enabled but discarded `make apiserver > /dev/null 2>&1`, all dependencies run under `docker-compose`
- database tuned on `docker-compose.yml` to be able handling high load
- **Write** = insert into 1 table, retrieve the id, append the id to global array (ignoring race condition)
- **Read** = read random id from global array, then query from 1 table
- **Health** = do syscall (or read from `/proc`) cached once per second
- **Hello** = only serializing empty input and rendering `{"hello":"world"}` output
- Benchmark server: 32-core, 128GB RAM, NVMe disk
- comparison: [TechEmpower](//www.techempower.com/benchmarks/#section=data-r20&hw=ph&test=update&l=zijocf-sf&a=2) which can do write 11K rps with spec: Intel Xeon Gold 5120 (28-core), 32 GB RAM, enterprise SSD for similar write use case.

## C10 Write (database write, network call)

```shell
 hey -n 100000 -c 10 http://localhost:9090/api/LoadTestWrite

Summary:
  Total:        2.1415 secs
  Slowest:      0.0276 secs
  Fastest:      0.0001 secs
  Average:      0.0002 secs
  Requests/sec: 46695.7216
  
  Total data:   6388900 bytes
  Size/request: 63 bytes

Response time histogram:
  0.000 [1]     |
  0.003 [99989] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.006 [0]     |
  0.008 [0]     |
  0.011 [0]     |
  0.014 [0]     |
  0.017 [0]     |
  0.019 [0]     |
  0.022 [0]     |
  0.025 [0]     |
  0.028 [10]    |


Latency distribution:
  10% in 0.0002 secs
  25% in 0.0002 secs
  50% in 0.0002 secs
  75% in 0.0002 secs
  90% in 0.0003 secs
  95% in 0.0003 secs
  99% in 0.0004 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.0276 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0003 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0004 secs
  resp wait:    0.0002 secs, 0.0001 secs, 0.0275 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0005 secs

Status code distribution:
  [200] 100000 responses
```

## C1K Write

```shell
hey -n 100000 -c 1000 http://localhost:9090/api/LoadTestWrite

Summary:
  Total:        1.7947 secs
  Slowest:      0.1562 secs
  Fastest:      0.0002 secs
  Average:      0.0174 secs
  Requests/sec: 55719.9159
  
  Total data:   6500000 bytes
  Size/request: 65 bytes

Response time histogram:
  0.000 [1]     |
  0.016 [78824] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.031 [4794]  |■■
  0.047 [692]   |
  0.063 [13757] |■■■■■■■
  0.078 [1569]  |■
  0.094 [91]    |
  0.109 [105]   |
  0.125 [7]     |
  0.141 [105]   |
  0.156 [55]    |


Latency distribution:
  10% in 0.0062 secs
  25% in 0.0077 secs
  50% in 0.0098 secs
  75% in 0.0139 secs
  90% in 0.0556 secs
  95% in 0.0584 secs
  99% in 0.0656 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0002 secs, 0.0002 secs, 0.1562 secs
  DNS-lookup:   0.0003 secs, 0.0000 secs, 0.1368 secs
  req write:    0.0001 secs, 0.0000 secs, 0.1368 secs
  resp wait:    0.0148 secs, 0.0001 secs, 0.0678 secs
  resp read:    0.0013 secs, 0.0000 secs, 0.1377 secs

Status code distribution:
  [200] 100000 responses
```

## C 10K Write

```shell
hey -n 100000 -c 10000 http://localhost:9090/api/LoadTestWrite

Summary:
  Total:        1.9145 secs
  Slowest:      0.7293 secs
  Fastest:      0.0002 secs
  Average:      0.1690 secs
  Requests/sec: 52232.5421
  
  Total data:   6500000 bytes
  Size/request: 65 bytes

Response time histogram:
  0.000 [1]     |
  0.073 [6625]  |■■■■■
  0.146 [29891] |■■■■■■■■■■■■■■■■■■■■■■■■■
  0.219 [48295] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.292 [9863]  |■■■■■■■■
  0.365 [3956]  |■■■
  0.438 [972]   |■
  0.511 [290]   |
  0.583 [71]    |
  0.656 [28]    |
  0.729 [8]     |


Latency distribution:
  10% in 0.0925 secs
  25% in 0.1315 secs
  50% in 0.1746 secs
  75% in 0.1927 secs
  90% in 0.2576 secs
  95% in 0.2950 secs
  99% in 0.3889 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0097 secs, 0.0002 secs, 0.7293 secs
  DNS-lookup:   0.0107 secs, 0.0000 secs, 0.1968 secs
  req write:    0.0014 secs, 0.0000 secs, 0.1788 secs
  resp wait:    0.1175 secs, 0.0001 secs, 0.2241 secs
  resp read:    0.0214 secs, 0.0000 secs, 0.4949 secs

Status code distribution:
  [200] 100000 responses
```

## C 20K Write

```shell
hey -n 100000 -c 20000 http://localhost:9090/api/LoadTestWrite

Summary:
  Total:        2.9915 secs
  Slowest:      2.5130 secs
  Fastest:      0.0003 secs
  Average:      0.5248 secs
  Requests/sec: 33428.1736
  
  Total data:   6500000 bytes
  Size/request: 65 bytes

Response time histogram:
  0.000 [1]     |
  0.252 [45724] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.503 [22434] |■■■■■■■■■■■■■■■■■■■■
  0.754 [8519]  |■■■■■■■
  1.005 [3409]  |■■■
  1.257 [2587]  |■■
  1.508 [6568]  |■■■■■■
  1.759 [8467]  |■■■■■■■
  2.010 [1971]  |■■
  2.262 [252]   |
  2.513 [68]    |


Latency distribution:
  10% in 0.0490 secs
  25% in 0.1373 secs
  50% in 0.3183 secs
  75% in 0.6613 secs
  90% in 1.5236 secs
  95% in 1.6808 secs
  99% in 1.9080 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0707 secs, 0.0003 secs, 2.5130 secs
  DNS-lookup:   0.0454 secs, 0.0000 secs, 0.9131 secs
  req write:    0.0058 secs, 0.0000 secs, 0.4495 secs
  resp wait:    0.2291 secs, 0.0002 secs, 0.9176 secs
  resp read:    0.0951 secs, 0.0000 secs, 1.0243 secs

Status code distribution:
  [200] 100000 responses
```

## C 10 Read (database query, network call)

```shell
hey -n 100000 -c 10 http://localhost:9090/api/LoadTestRead 

Summary:
  Total:        1.8199 secs
  Slowest:      0.0082 secs
  Fastest:      0.0001 secs
  Average:      0.0002 secs
  Requests/sec: 54946.8908
  
  Total data:   36200000 bytes
  Size/request: 362 bytes

Response time histogram:
  0.000 [1]     |
  0.001 [99924] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.002 [24]    |
  0.003 [12]    |
  0.003 [9]     |
  0.004 [8]     |
  0.005 [1]     |
  0.006 [11]    |
  0.007 [0]     |
  0.007 [0]     |
  0.008 [10]    |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0002 secs
  50% in 0.0002 secs
  75% in 0.0002 secs
  90% in 0.0002 secs
  95% in 0.0002 secs
  99% in 0.0004 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0001 secs, 0.0082 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0003 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0004 secs
  resp wait:    0.0002 secs, 0.0001 secs, 0.0082 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0005 secs

Status code distribution:
  [200] 100000 responses
```

## C1K Read

```shell
Summary:
  Total:        0.8863 secs
  Slowest:      0.2039 secs
  Fastest:      0.0001 secs
  Average:      0.0082 secs
  Requests/sec: 112834.5020
  
  Total data:   36200000 bytes
  Size/request: 362 bytes

Response time histogram:
  0.000 [1]     |
  0.020 [89162] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.041 [9143]  |■■■■
  0.061 [650]   |
  0.082 [62]    |
  0.102 [67]    |
  0.122 [515]   |
  0.143 [161]   |
  0.163 [68]    |
  0.183 [165]   |
  0.204 [6]     |


Latency distribution:
  10% in 0.0003 secs
  25% in 0.0004 secs
  50% in 0.0008 secs
  75% in 0.0113 secs
  90% in 0.0211 secs
  95% in 0.0261 secs
  99% in 0.0761 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0006 secs, 0.0001 secs, 0.2039 secs
  DNS-lookup:   0.0004 secs, 0.0000 secs, 0.1535 secs
  req write:    0.0003 secs, 0.0000 secs, 0.1522 secs
  resp wait:    0.0007 secs, 0.0001 secs, 0.0604 secs
  resp read:    0.0036 secs, 0.0000 secs, 0.0811 secs

Status code distribution:
  [200] 100000 responses
```

## C 10K Read

```shell
hey -n 100000 -c 10000 http://localhost:9090/api/LoadTestRead 

Summary:
  Total:        1.2142 secs
  Slowest:      1.1184 secs
  Fastest:      0.0001 secs
  Average:      0.1015 secs
  Requests/sec: 82359.4925
  
  Total data:   36200000 bytes
  Size/request: 362 bytes

Response time histogram:
  0.000 [1]     |
  0.112 [63986] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.224 [19104] |■■■■■■■■■■■■
  0.336 [8328]  |■■■■■
  0.447 [4517]  |■■■
  0.559 [3001]  |■■
  0.671 [708]   |
  0.783 [278]   |
  0.895 [69]    |
  1.007 [6]     |
  1.118 [2]     |


Latency distribution:
  10% in 0.0003 secs
  25% in 0.0004 secs
  50% in 0.0050 secs
  75% in 0.1653 secs
  90% in 0.3115 secs
  95% in 0.4017 secs
  99% in 0.5704 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0105 secs, 0.0001 secs, 1.1184 secs
  DNS-lookup:   0.0066 secs, 0.0000 secs, 0.3608 secs
  req write:    0.0035 secs, 0.0000 secs, 0.2850 secs
  resp wait:    0.0129 secs, 0.0001 secs, 0.2628 secs
  resp read:    0.0378 secs, 0.0000 secs, 0.7637 secs

Status code distribution:
  [200] 100000 responses
```

## C 20K Read

```shell
hey -n 100000 -c 20000 http://localhost:9090/api/LoadTestRead 

Summary:
  Total:        2.5815 secs
  Slowest:      2.3408 secs
  Fastest:      0.0002 secs
  Average:      0.4402 secs
  Requests/sec: 38736.6084
  
  Total data:   36200000 bytes
  Size/request: 362 bytes

Response time histogram:
  0.000 [1]     |
  0.234 [47639] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.468 [16885] |■■■■■■■■■■■■■■
  0.702 [6547]  |■■■■■
  0.936 [7330]  |■■■■■■
  1.171 [11008] |■■■■■■■■■
  1.405 [3024]  |■■■
  1.639 [4877]  |■■■■
  1.873 [2055]  |■■
  2.107 [460]   |
  2.341 [174]   |


Latency distribution:
  10% in 0.0004 secs
  25% in 0.0005 secs
  50% in 0.2490 secs
  75% in 0.7878 secs
  90% in 1.1858 secs
  95% in 1.5508 secs
  99% in 1.8054 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0483 secs, 0.0002 secs, 2.3408 secs
  DNS-lookup:   0.0620 secs, 0.0000 secs, 0.8380 secs
  req write:    0.0032 secs, 0.0000 secs, 0.5570 secs
  resp wait:    0.1060 secs, 0.0001 secs, 0.7873 secs
  resp read:    0.1263 secs, 0.0000 secs, 1.2577 secs

Status code distribution:
  [200] 100000 responses
```

## C 10 Health (no database, syscall once per second)

```shell
 hey -n 100000 -c 10 http://localhost:9090/api/Health      

Summary:
  Total:        0.9407 secs
  Slowest:      0.0089 secs
  Fastest:      0.0000 secs
  Average:      0.0001 secs
  Requests/sec: 106305.0847
  
  Total data:   9897468 bytes
  Size/request: 98 bytes

Response time histogram:
  0.000 [1]     |
  0.001 [99950] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.002 [26]    |
  0.003 [1]     |
  0.004 [1]     |
  0.004 [1]     |
  0.005 [0]     |
  0.006 [2]     |
  0.007 [8]     |
  0.008 [0]     |
  0.009 [10]    |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0001 secs
  50% in 0.0001 secs
  75% in 0.0001 secs
  90% in 0.0001 secs
  95% in 0.0001 secs
  99% in 0.0003 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0089 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0002 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0017 secs
  resp wait:    0.0001 secs, 0.0000 secs, 0.0089 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0011 secs

Status code distribution:
  [200] 100000 responses
```

## C1K Health

```shell
hey -n 100000 -c 1000 http://localhost:9090/api/Health 

Summary:
  Total:        0.6982 secs
  Slowest:      0.1649 secs
  Fastest:      0.0001 secs
  Average:      0.0065 secs
  Requests/sec: 143231.8962
  
  Total data:   9900000 bytes
  Size/request: 99 bytes

Response time histogram:
  0.000 [1]     |
  0.017 [87831] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.033 [9472]  |■■■■
  0.050 [1401]  |■
  0.066 [426]   |
  0.082 [470]   |
  0.099 [284]   |
  0.115 [15]    |
  0.132 [72]    |
  0.148 [23]    |
  0.165 [5]     |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0001 secs
  50% in 0.0002 secs
  75% in 0.0092 secs
  90% in 0.0181 secs
  95% in 0.0236 secs
  99% in 0.0600 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0004 secs, 0.0001 secs, 0.1649 secs
  DNS-lookup:   0.0004 secs, 0.0000 secs, 0.0744 secs
  req write:    0.0002 secs, 0.0000 secs, 0.0696 secs
  resp wait:    0.0002 secs, 0.0000 secs, 0.0404 secs
  resp read:    0.0028 secs, 0.0000 secs, 0.0876 secs

Status code distribution:
  [200] 100000 responses
```

## C 10K Health

```shell
hey -n 100000 -c 10000 http://localhost:9090/api/Health

Summary:
  Total:        1.0075 secs
  Slowest:      0.8499 secs
  Fastest:      0.0001 secs
  Average:      0.0844 secs
  Requests/sec: 99251.1317
  
  Total data:   9870145 bytes
  Size/request: 98 bytes

Response time histogram:
  0.000 [1]     |
  0.085 [65039] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.170 [17363] |■■■■■■■■■■■
  0.255 [6213]  |■■■■
  0.340 [6799]  |■■■■
  0.425 [2234]  |■
  0.510 [1238]  |■
  0.595 [716]   |
  0.680 [287]   |
  0.765 [88]    |
  0.850 [22]    |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0001 secs
  50% in 0.0029 secs
  75% in 0.1364 secs
  90% in 0.2713 secs
  95% in 0.3344 secs
  99% in 0.5148 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0085 secs, 0.0001 secs, 0.8499 secs
  DNS-lookup:   0.0053 secs, 0.0000 secs, 0.1781 secs
  req write:    0.0009 secs, 0.0000 secs, 0.1783 secs
  resp wait:    0.0039 secs, 0.0000 secs, 0.1745 secs
  resp read:    0.0329 secs, 0.0000 secs, 0.5704 secs

Status code distribution:
  [200] 100000 responses
```

## C 20K Health

```shell
hey -n 100000 -c 20000 http://localhost:9090/api/Health      

Summary:
  Total:        1.9356 secs
  Slowest:      1.8303 secs
  Fastest:      0.0001 secs
  Average:      0.2596 secs
  Requests/sec: 51664.5346
  
  Total data:   9888826 bytes
  Size/request: 98 bytes

Response time histogram:
  0.000 [1]     |
  0.183 [57872] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.366 [17025] |■■■■■■■■■■■■
  0.549 [7560]  |■■■■■
  0.732 [5768]  |■■■■
  0.915 [4553]  |■■■
  1.098 [4404]  |■■■
  1.281 [1533]  |■
  1.464 [434]   |
  1.647 [735]   |■
  1.830 [115]   |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0002 secs
  50% in 0.1550 secs
  75% in 0.3710 secs
  90% in 0.7834 secs
  95% in 0.9897 secs
  99% in 1.4058 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0070 secs, 0.0001 secs, 1.8303 secs
  DNS-lookup:   0.0229 secs, 0.0000 secs, 0.6348 secs
  req write:    0.0017 secs, 0.0000 secs, 0.4559 secs
  resp wait:    0.0376 secs, 0.0000 secs, 0.5288 secs
  resp read:    0.0887 secs, 0.0000 secs, 1.4141 secs

Status code distribution:
  [200] 100000 responses
```

## C 10 Hello (serialization-deserialization only)

```shell
hey -n 100000 -c 10 http://localhost:9090/api/LoadHello

Summary:
  Total:        0.9550 secs
  Slowest:      0.0085 secs
  Fastest:      0.0000 secs
  Average:      0.0001 secs
  Requests/sec: 104708.9702
  
  Total data:   5700000 bytes
  Size/request: 57 bytes

Response time histogram:
  0.000 [1]     |
  0.001 [99911] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.002 [45]    |
  0.003 [4]     |
  0.003 [8]     |
  0.004 [1]     |
  0.005 [10]    |
  0.006 [10]    |
  0.007 [0]     |
  0.008 [0]     |
  0.008 [10]    |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0001 secs
  50% in 0.0001 secs
  75% in 0.0001 secs
  90% in 0.0001 secs
  95% in 0.0001 secs
  99% in 0.0003 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0085 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0003 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0004 secs
  resp wait:    0.0001 secs, 0.0000 secs, 0.0084 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0007 secs

Status code distribution:
  [200] 100000 responses
```

## C1K Hello

```shell
hey -n 100000 -c 1000 http://localhost:9090/api/LoadHello 

Summary:
  Total:        0.6930 secs
  Slowest:      0.1352 secs
  Fastest:      0.0001 secs
  Average:      0.0063 secs
  Requests/sec: 144292.8724
  
  Total data:   5700000 bytes
  Size/request: 57 bytes

Response time histogram:
  0.000 [1]     |
  0.014 [85481] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.027 [10190] |■■■■■
  0.041 [3636]  |■■
  0.054 [549]   |
  0.068 [83]    |
  0.081 [30]    |
  0.095 [7]     |
  0.108 [6]     |
  0.122 [7]     |
  0.135 [10]    |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0001 secs
  50% in 0.0002 secs
  75% in 0.0093 secs
  90% in 0.0193 secs
  95% in 0.0244 secs
  99% in 0.0389 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0001 secs, 0.0001 secs, 0.1352 secs
  DNS-lookup:   0.0003 secs, 0.0000 secs, 0.0492 secs
  req write:    0.0001 secs, 0.0000 secs, 0.0755 secs
  resp wait:    0.0003 secs, 0.0000 secs, 0.0631 secs
  resp read:    0.0027 secs, 0.0000 secs, 0.0839 secs

Status code distribution:
  [200] 100000 responses
```

## C 10K Hello

```shell
hey -n 100000 -c 10000 http://localhost:9090/api/LoadHello

Summary:
  Total:        1.0122 secs
  Slowest:      0.9236 secs
  Fastest:      0.0001 secs
  Average:      0.0828 secs
  Requests/sec: 98795.5020
  
  Total data:   5700000 bytes
  Size/request: 57 bytes

Response time histogram:
  0.000 [1]     |
  0.092 [64914] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.185 [18645] |■■■■■■■■■■■
  0.277 [6256]  |■■■■
  0.369 [5515]  |■■■
  0.462 [3459]  |■■
  0.554 [810]   |
  0.647 [318]   |
  0.739 [55]    |
  0.831 [24]    |
  0.924 [3]     |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0001 secs
  50% in 0.0004 secs
  75% in 0.1227 secs
  90% in 0.2803 secs
  95% in 0.3552 secs
  99% in 0.4940 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0116 secs, 0.0001 secs, 0.9236 secs
  DNS-lookup:   0.0095 secs, 0.0000 secs, 0.1775 secs
  req write:    0.0010 secs, 0.0000 secs, 0.1084 secs
  resp wait:    0.0094 secs, 0.0000 secs, 0.1262 secs
  resp read:    0.0272 secs, 0.0000 secs, 0.5055 secs

Status code distribution:
  [200] 100000 responses
```

## C 20K Hello

```shell
hey -n 100000 -c 20000 http://localhost:9090/api/LoadHello

Summary:
  Total:        2.2549 secs
  Slowest:      2.2309 secs
  Fastest:      0.0001 secs
  Average:      0.3807 secs
  Requests/sec: 44347.0950
  
  Total data:   5700000 bytes
  Size/request: 57 bytes

Response time histogram:
  0.000 [1]     |
  0.223 [48104] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.446 [15563] |■■■■■■■■■■■■■
  0.669 [9233]  |■■■■■■■■
  0.892 [13934] |■■■■■■■■■■■■
  1.115 [7301]  |■■■■■■
  1.339 [707]   |■
  1.562 [3671]  |■■■
  1.785 [1016]  |■
  2.008 [405]   |
  2.231 [65]    |


Latency distribution:
  10% in 0.0001 secs
  25% in 0.0002 secs
  50% in 0.2425 secs
  75% in 0.7113 secs
  90% in 0.9720 secs
  95% in 1.3476 secs
  99% in 1.6238 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0351 secs, 0.0001 secs, 2.2309 secs
  DNS-lookup:   0.0494 secs, 0.0000 secs, 1.3221 secs
  req write:    0.0091 secs, 0.0000 secs, 1.0711 secs
  resp wait:    0.0723 secs, 0.0000 secs, 1.0124 secs
  resp read:    0.1177 secs, 0.0000 secs, 1.3883 secs

Status code distribution:
  [200] 100000 responses
```
