## UI

You can [watch the video here](https://drive.google.com/file/d/1nPJpz9YAmwzWLcWc5Qqk90djmZsTZv3w/view?usp=sharing).

## Order Service

### Quick start
````bash
git clone https://github.com/Parside01/wb_tech_l0.git
cd wb_tech_l0
````
````bash
make app-run-docker
````

### More options

#### Running tests
````bash
make run-tests
````

#### Coverage
````bash
make coverage 
````

#### Run linter
````bash
make lint
````

## Benchmarks
#### Before starting, we will add `uuids` in `benchmarks/uuids.txt`
#### WRK

````bash
wrk -t12 -c400 -d30s -s benchmarks/wrk/wrk-script.lua http://localhost:8080
````
````bash
12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    17.06ms   13.72ms 147.46ms   78.02%
    Req/Sec     1.37k   802.29     3.41k    67.50%
  24344 requests in 30.90s, 37.27MB read
  Socket errors: connect 0, read 0, write 0, timeout 792
  Non-2xx or 3xx responses: 12
Requests/sec:    787.94
Transfer/sec:      1.21MB
````

#### Vegeta
````bash
./benchmarks/vegeta/generate_req.sh
cat benchmarks/vegeta/req.txt | vegeta attack -rate=1000 -duration=30s | vegeta report 
````
````bash
Requests      [total, rate, throughput]         29999, 1000.00, 570.35
Duration      [total, attack, wait]             52.545s, 29.999s, 22.546s
Latencies     [min, mean, 50, 90, 95, 99, max]  403.024µs, 1.576s, 1.033ms, 7.261s, 8.669s, 10.368s, 22.699s
Bytes In      [total, mean]                     44529739, 1484.37
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           99.90%
Status Codes  [code:count]                      200:29969  500:30
Error Set:
````
