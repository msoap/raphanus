Raphanus - simple in-memory cache
---------------------------------

[![GoDoc](https://godoc.org/github.com/msoap/raphanus?status.svg)](https://godoc.org/github.com/msoap/raphanus)
[![Build Status](https://travis-ci.org/msoap/raphanus.svg?branch=master)](https://travis-ci.org/msoap/raphanus)
[![Coverage Status](https://coveralls.io/repos/github/msoap/raphanus/badge.svg?branch=master)](https://coveralls.io/github/msoap/raphanus?branch=master)
[![Docker Pulls](https://img.shields.io/docker/pulls/msoap/raphanus.svg?maxAge=3600)](https://hub.docker.com/r/msoap/raphanus/)
[![Sourcegraph](https://sourcegraph.com/github.com/msoap/raphanus/-/badge.svg)](https://sourcegraph.com/github.com/msoap/raphanus?badge)
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

in Docker container:

    docker run --name raphanus --publish 8771:8771 --detach msoap/raphanus

Mirror docker [repository](https://quay.io/repository/msoap/raphanus): `quay.io/msoap/raphanus:latest`

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
    
    Benchmark_raphanusServer-4      	    5000	   1196515 ns/op	    9432 B/op	     128 allocs/op
    Benchmark_raphanusEmbed-4       	 5000000	      1633 ns/op	     184 B/op	       4 allocs/op
    Benchmark_redis-4               	   10000	    636079 ns/op	     343 B/op	      19 allocs/op
    Benchmark_memcache-4            	   10000	    689439 ns/op	    2744 B/op	      63 allocs/op
    Benchmark_raphanusServerTTL-4   	    5000	   1178499 ns/op	    9459 B/op	     130 allocs/op
    Benchmark_raphanusEmbedTTL-4    	 5000000	      2554 ns/op	     194 B/op	       6 allocs/op
    Benchmark_redisTTL-4            	   10000	    791056 ns/op	     417 B/op	      21 allocs/op
    Benchmark_memcacheTTL-4         	   10000	    733611 ns/op	    2744 B/op	      63 allocs/op

### local raphanus and redis servers (on MacOS)

Redis 3.2.7

memcached 1.4.34

    $ make run-benchmark
    Benchmark_raphanusServer-4      	   20000	    349394 ns/op	    9465 B/op	     131 allocs/op
    Benchmark_raphanusEmbed-4       	 3000000	      1919 ns/op	     213 B/op	       7 allocs/op
    Benchmark_boltdb-4              	   30000	    258030 ns/op	   34320 B/op	      69 allocs/op
    Benchmark_redis-4               	   50000	    128278 ns/op	     411 B/op	      23 allocs/op
    Benchmark_memcache-4            	   50000	    135533 ns/op	    2752 B/op	      63 allocs/op
    Benchmark_raphanusServerTTL-4   	   20000	    356989 ns/op	    9499 B/op	     133 allocs/op
    Benchmark_raphanusEmbedTTL-4    	 2000000	      3347 ns/op	     244 B/op	       9 allocs/op
    Benchmark_redisTTL-4            	   50000	    127851 ns/op	     473 B/op	      25 allocs/op
    Benchmark_memcacheTTL-4         	   50000	    135326 ns/op	    2811 B/op	      67 allocs/op
