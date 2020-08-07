# Hash Service

Password Hashing Service

## Description

A very simple password hashing service with concurrent support.

## Installing

### Local

The hash service has only one dependency (`github.com/stretchr/testify`) which is used for tests only. 
To build the service locally, simply run:

```bash
$ make service
go build -o <GOPATH>/src/github.com/plar/hash/_bin/hashsvc ./cmd/hashsvc/main.go
```

Now you can use the `./_bin/hashsvc` to start the service.

To run tests (with the race detector), run the following command:
```bash
$ make test
?   	github.com/plar/hash/cmd/hashsvc	[no test files]
ok  	github.com/plar/hash/domain	0.032s	coverage: 100.0% of statements
?   	github.com/plar/hash/domain/repository	[no test files]
ok  	github.com/plar/hash/infra/persistence/memory	0.033s	coverage: 100.0% of statements
ok  	github.com/plar/hash/infra/pool	0.808s	coverage: 84.1% of statements
?   	github.com/plar/hash/server	[no test files]
?   	github.com/plar/hash/server/config	[no test files]
ok  	github.com/plar/hash/service/hasher	0.035s	coverage: 53.8% of statements
ok  	github.com/plar/hash/service/health	0.041s	coverage: 100.0% of statements
?   	github.com/plar/hash/service/stats	[no test files]
```

### Docker

You also can use `docker` to build and run the service.

```
$ make docker-service
Sending build context to Docker daemon  8.316MB
- - 8< - - - 8< - - - 8< - - - 8< - - - 8< - - - 8< - - - 8< - - - 8< --
Successfully built b7dfd34662d6
Successfully tagged plar/hashsvc:latest
```

Now you can use the `plar/hashsvc:latest` image to start the service.

## Service Parameters

## Environment variables

| Parameter | Description | Valid Values | Default |
|-----------|-------------|--------------|---------|
|HASH_SERVER_ADDR| the address for the service to listen on | IP or Host name | 0.0.0.0 |
|HASH_SERVER_PORT| the service port to listen on | 1..65535 | 8080 |
|HASH_SERVER_READ_TIMEOUT| the maximum duration in seconds for reading the entire request, including the body | Integers | 10 |
|HASH_SERVER_WRITE_TIMEOUT| the maximum duration in seconds before timing out writes of the response | Integers| 10 |
|HASH_SERVER_IDLE_TIMEOUT | the maximum amount of time in seconds to wait for the next request when keep-alives are enabled | Integers | 10 |
|HASH_SERVER_SHUTDOWN_TIMEOUT| the maximum duration in seconds to complete a graceful shutdown | Integers | 30 |
|HASH_TASK_DELAY| number of seconds to delay hash task | Integers | 5 |
|HASH_TOTAL_WORKERS| number of workers in the hash pool | Positive integers | Number of logical CPU cores on the running system |
|HASH_QUEUE_SIZE| hash pool queue size | Positive integers | 10000 |

## CLI Arguments

| Parameter | Description | Valid Values | Default |
|-----------|-------------|--------------|---------|
| `-delay`       | number of seconds to delay hash task | Integers | 5 |
| `-workers` | number of workers in the hash pool | Positive integers | Number of logical CPU cores on the running system |
| `-queue-size` | hash pool queue size | Positive integers | 10000 |


## Endpoints
| Method | Endpoint | Request Content-Type | URI Parameters | Request | Response |
|--------|----------|----------------------|----------------|---------|----------|
| `POST` | `/hash`     | application/x-www-form-urlencoded | - | A `password` for hashing.<br> Example: `angryMonkey` | An integer number of job ID.<br> Example: `1` |
| `GET`  | `/hash/{id}`| text/plain | `id` the integer job ID | - | If found, base64 encoded string of the SHA512 hash for the job ID. <br> Example: `YW5ncnlNb25rZXnPg+...` |
| `GET`  | `/stats`    | text/plain | - | - | A basic information about your password hashes.<br> Eaxample: `{"total":149,"average":2}` |
| `GET`  | `/shutdown` | text/plain | - | - | Graceful shutdown request. The service will wait for any pending/in-flight work to finish before exiting and will reject new requests as well.

## Examples

Lets run the service with a 30 second task delay, request some hashes, check the resulting hashes, check stats, and shutdown the service.

### Run the service

```bash
$ ./_bin/hashsvc -delay 30
```

or if you built the service using `docker` then run the following command:

```bash
$ docker run -p 8080:8080 plar/hashsvc:latest -delay 30
```


And after you will see something like this:
```
2020/08/07 12:24:30 The service is ready to handle requests at 0.0.0.0:8080
```

### Request some hashes

Now we are ready to send requests.

```bash
$ curl --data "password=angryMonkey" http://localhost:8080/hash
1
```

You will see a job ID (`1`) as a response.

### Check the resulting hashes

Let's try to request the resulting hash by ID.

```bash
$ curl http://localhost:8080/hash/1
```

If 30 seconds have not passed since the hashing request or you passed a non-existent identifier, the service will return a 404 error.

```
HashID '1': hash not found
```

Let's try again in 30 seconds.

```bash
$ curl http://localhost:8080/hash/1
YW5ncnlNb25rZXnPg+E1fu+4vfFUKFDWbYAH1iDkBQtXFdyD9Kkh02zpzkfQ0TxdhfKw/4MY0od+7C9juTG9R0F6gaU4Mnr5J9o+
```

Now we can see that the service returned us a base64/SHA512 encoded hash for the `angryMonkey` password.

### Check the service stats
```bash
$ curl http://localhost:8080/stats
{"total":1,"average":2}
```

We have sent one request only and and it took 2 microseconds.

### Shutdown the service

```bash
$ curl http://localhost:8080/shutdown
```

You should see the similar log lines in the service log

```
2020/08/07 12:43:39 the stats service method=TrackMetric name=Hasher.Get, count=1, execTime=36.12µs
2020/08/07 12:43:39 the hasher service method=Get id=1 => hash=hash{}, err=HashID '1': hash not found
2020/08/07 12:43:44 the stats service method=TrackMetric name=Hasher.Create, count=1, execTime=2.286µs
2020/08/07 12:43:44 the hasher service method=Create => id=1, err=<nil>
2020/08/07 12:43:48 the stats service method=Metric name=Hasher.Create => {1 2286}
2020/08/07 12:52:07 the server is shutting down
2020/08/07 12:52:07 the hasher service is stopping
2020/08/07 12:52:36 pool worker #4: quit, processed tasks=1, total exec time=30.000336537s
2020/08/07 12:52:36 the hasher service has been stopped
2020/08/07 12:52:36 pool worker #7: quit, processed tasks=0, total exec time=0s
2020/08/07 12:52:36 pool worker #5: quit, processed tasks=0, total exec time=0s
2020/08/07 12:52:36 pool worker #3: quit, processed tasks=0, total exec time=0s
2020/08/07 12:52:36 pool worker #6: quit, processed tasks=0, total exec time=0s
2020/08/07 12:52:36 pool worker #2: quit, processed tasks=0, total exec time=0s
2020/08/07 12:52:36 pool worker #1: quit, processed tasks=0, total exec time=0s
2020/08/07 12:52:36 pool worker #8: quit, processed tasks=0, total exec time=0s
2020/08/07 12:52:36 The server has been shutdown
```

Let's try to send a new hash request again while our service is trying to shutdown.
```
$ curl --data "password=happyMonkey" http://localhost:8080/hash
the server is shutting down
```

Note that approximately `delay`(30) seconds  elapsed between 2 events:

```
2020/08/07 12:52:07 the server is shutting down
...
2020/08/07 12:52:36 The server has been shutdown
```
