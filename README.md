Raphanus - simple Redis-like in-memory cache
--------------------------------------------
[![GoDoc](https://godoc.org/github.com/msoap/raphanus?status.svg)](https://godoc.org/github.com/msoap/raphanus)
[![Build Status](https://travis-ci.org/msoap/raphanus.svg?branch=master)](https://travis-ci.org/msoap/raphanus)
[![Coverage Status](https://coveralls.io/repos/github/msoap/raphanus/badge.svg?branch=master)](https://coveralls.io/github/msoap/raphanus?branch=master)
[![Docker Pulls](https://img.shields.io/docker/pulls/msoap/raphanus.svg?maxAge=3600)](https://hub.docker.com/r/msoap/raphanus/)
[![Report Card](https://goreportcard.com/badge/github.com/msoap/raphanus)](https://goreportcard.com/report/github.com/msoap/raphanus)

## Install

From source:

    go get -u github.com/msoap/raphanus
    # build server & cli
    cd $GOPATH/src/github.com/msoap/raphanus/server && go build -o $GOPATH/bin/raphanus-server
    # cli client not implement now
    cd $GOPATH/src/github.com/msoap/raphanus/cli && go build -o $GOPATH/bin/raphanus-cli

From Docker hub:

    docker pull msoap/raphanus

## Run server

    raphanus-server [options]
    options:
      -address string
           	address for bind server (default "localhost:8771")
      -auth string
           	user:password for enable HTTP basic authentication
      -filename string
           	file name for storage on disk, '' - for work in-memory only
      -sync-time int
           	time in seconds between sync on disk
      -version
           	get version

as Docker container:

    docker run --name raphanus --publish 8771:8771 --detach msoap/raphanus

## Examples: get calls to server by curl

 * get count of keys: `curl -s 'http://localhost:8771/v1/length'`
 * get keys: `curl -s 'http://localhost:8771/v1/keys'`
 * get stat with authentication: `curl -s -u user:pass 'http://localhost:8771/v1/stat'`
 * set integer key `k1` with ttl (100 sec): `curl -s -X POST -d 123 'http://localhost:8771/v1/int/k1?ttl=100'`
 * get integer key `k1`: `curl -s 'http://localhost:8771/v1/int/k1'`
 * set string key `k2` without ttl: `curl -s -X POST -d 'str value' 'http://localhost:8771/v1/str/k2'`
 * get string key `k2`: `curl -s 'http://localhost:8771/v1/str/k2'`
 * set list value: `curl -s -X POST -H 'Content-Type: application/json' -d '["v1", "v2"]' http://localhost:8771/v1/list/k3`
 * get list value: `curl -s 'http://localhost:8771/v1/list/k3'`
 * set dict value: `curl -s -X POST -H 'Content-Type: application/json' -d '{"dk1": "v1", "dk2": "v2"}' http://localhost:8771/v1/dict/k4`
 * get dict value: `curl -s 'http://localhost:8771/v1/dict/k4'`
 * delete key `k1`: `curl -s -X DELETE http://localhost:8771/v1/remove/k1`
 * see other in [handlers.go](https://github.com/msoap/raphanus/blob/master/server/handlers.go)

## Use as embed DB

```Go
import (
    "github.com/msoap/raphanus"
    "github.com/msoap/raphanus/common"
)

func main() {
    raph := raphanus.New("", 0)
    // or with storage, with sync every 300 seconds
    // raph := raphanus.New("filename.db", 300)
    raph.SetStr("key", "value")
    v, err := raph.GetStr("key")
    if err == raphanuscommon.ErrKeyNotExists {
        ...
    }

    raph.UnderLock(func () {
        v, err := raph.GetStr("k1")
        if err != nil {
            return
        }
        raph.SetStr("k1", v + " updated", 0)
    })
}
```

## Use client library for connect with server
[![GoDoc](https://godoc.org/github.com/msoap/raphanus/client?status.svg)](https://godoc.org/github.com/msoap/raphanus/client)

```Go
import (
    "github.com/msoap/raphanus/client"
    "github.com/msoap/raphanus/common"
)

func main() {
	// with default address:
	raph := raphanusclient.New()
	// or with another address:
	// raph := raphanusclient.New(raphanusclient.Cfg{Address: "http://localhost:8771"})
	// or with authentication:
	// raph := raphanusclient.New(raphanusclient.Cfg{User: "uname", Password: "pass"})

    raph.SetStr("key", "value", 3600)
    v, err := raph.GetStr("key")
    if err == raphanuscommon.ErrKeyNotExists {
        ...
    }
}
```

An example of using the library: [simple.go](https://github.com/msoap/raphanus/blob/master/client/examples/simple.go)

## Bencmarks:
### with servers in Docker

    $ docker run --name raphanus --publish 8771:8771 --rm msoap/raphanus
    $ docker run --name redis --rm --publish 6379:6379 redis
    $ docker run --name memcache --rm --publish 11211:11211 memcached
    $ make run-benchmark
    
    Benchmark_raphanusServer-4     	    2000       	    909746 ns/op
    Benchmark_raphanusEmbed-4      	 1000000       	      2307 ns/op
    Benchmark_redis-4              	    2000       	    592827 ns/op
    Benchmark_memcache-4           	    2000       	    593110 ns/op
    Benchmark_raphanusServerTTL-4  	    2000       	   1011409 ns/op
    Benchmark_raphanusEmbedTTL-4   	  200000       	      8499 ns/op
    Benchmark_redisTTL-4           	    2000       	    681730 ns/op
    Benchmark_memcacheTTL-4        	    2000       	    534125 ns/op

### local raphanus and redis servers (on MacOS)

    $ make run-benchmark
    Benchmark_raphanusServer-4     	    5000       	    301054 ns/op
    Benchmark_raphanusEmbed-4      	 1000000       	      2309 ns/op
    Benchmark_redis-4              	   10000       	    127818 ns/op
    Benchmark_memcache-4           	   10000       	    172703 ns/op
    Benchmark_raphanusServerTTL-4  	    5000       	    380316 ns/op
    Benchmark_raphanusEmbedTTL-4   	  200000       	      9202 ns/op
    Benchmark_redisTTL-4           	   10000       	    149906 ns/op
    Benchmark_memcacheTTL-4        	   10000       	    143990 ns/op
