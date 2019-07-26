# Goal

Demonstrate how to use the [kit-mw](https://github.com/jwenz723/kit-mw) middlewares.

# Running the code

> This service is almost a direct copy of [addsvc](https://github.com/go-kit/kit/tree/master/examples/addsvc).

Start up the service

```bash
go run cmd/main.go 
```

Send in a `concat` and `sum` request

```bash
$ curl -d '{"a":"2","b":"3"}' localhost:8081/concat
  {"v":"23"}
$ curl -d '{"a":2,"b":3}' localhost:8081/sum       
  {"v":5}

```

See the keyvals printed for each request/response (the keys start with `ConcatRequest`, `ConcatResponse`, `SumRequest.`, and `SumResponse.`):

```bash
level=info ts=2019-07-26T08:00:38.638877Z caller=eplogger.go:56 method=Concat transport_error=null took=7.985µs ConcatRequest.A=2 ConcatRequest.B=3 ConcatResponse.V=23 ConcatResponse.Err=null
level=info ts=2019-07-26T08:00:44.300104Z caller=eplogger.go:56 method=Sum transport_error=null took=4.071µs SumRequest.A=2 SumRequest.B=3 SumResponse.R=5 SumResponse.Err=null
```
