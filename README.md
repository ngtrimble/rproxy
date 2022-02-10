# rproxy

## Introduction

This program proxies requests as follows:

```
http://localhost:8080/svca => http://localhost:8888/
http://localhost:8080/svcb => http://localhost:9999/
```

## Build

```
go build
```

## Run

```
./rproxy
```

## Test with curl

```
curl -v http://localhost:8080/svca
curl -v http://localhost:8080/svcb
curl -v http://localhost:8080/
```
