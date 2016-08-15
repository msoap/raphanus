/*
run:
	docker run --name raphanus --publish 8771:8771 --rm msoap/raphanus
	docker run --name redis --publish --rm 6379:6379 redis
	go test -bench Benchmark
*/
package raphanus_test

import (
	"strconv"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/msoap/raphanus"
	"github.com/msoap/raphanus/client"
)

func Benchmark_raphanusServer(b *testing.B) {
	raph := raphanusclient.New(raphanusclient.Cfg{Address: "http://localhost:8771"})

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := raph.SetStr("key_"+strI, "bar_"+strI, 0)
		if err != nil {
			b.Fatal(err)
		}

		newVal, err := raph.GetStr("key_" + strI)
		if err != nil {
			b.Fatal(err)
		}
		if newVal != "bar_"+strI {
			b.Fatal("Set/get not equal")
		}
	}
}

func Benchmark_raphanusEmbed(b *testing.B) {
	raph := raphanus.New("", 0)

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := raph.SetStr("key_"+strI, "bar_"+strI, 0)
		if err != nil {
			b.Fatal(err)
		}

		newVal, err := raph.GetStr("key_" + strI)
		if err != nil {
			b.Fatal(err)
		}
		if newVal != "bar_"+strI {
			b.Fatal("Set/get not equal")
		}
	}
}

func Benchmark_redis(b *testing.B) {
	redisCli, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err = redisCli.Cmd("SET", "key_"+strI, "bar_"+strI).Err
		if err != nil {
			b.Fatal(err)
		}

		newVal, err := redisCli.Cmd("GET", "key_"+strI).Str()
		if err != nil {
			b.Fatal(err)
		}
		if newVal != "bar_"+strI {
			b.Fatal("Set/get not equal")
		}
	}
}

func Benchmark_memcache(b *testing.B) {
	mc := memcache.New("localhost:11211")

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := mc.Set(&memcache.Item{Key: "key_" + strI, Value: []byte("bar_" + strI)})
		if err != nil {
			b.Fatal(err)
		}

		item, err := mc.Get("key_" + strI)
		if err != nil {
			b.Fatal(err)
		}
		newVal := item.Value
		if string(newVal) != "bar_"+strI {
			b.Fatal("Set/get not equal")
		}
	}
}

func Benchmark_raphanusServerTTL(b *testing.B) {
	raph := raphanusclient.New(raphanusclient.Cfg{Address: "http://localhost:8771"})

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := raph.SetStr("key_"+strI, "bar_"+strI, 2)
		if err != nil {
			b.Fatal(err)
		}

		newVal, err := raph.GetStr("key_" + strI)
		if err != nil {
			b.Fatal(err)
		}
		if newVal != "bar_"+strI {
			b.Fatal("Set/get not equal")
		}
	}
}

func Benchmark_raphanusEmbedTTL(b *testing.B) {
	raph := raphanus.New("", 0)

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := raph.SetStr("key_"+strI, "bar_"+strI, 2)
		if err != nil {
			b.Fatal(err)
		}

		newVal, err := raph.GetStr("key_" + strI)
		if err != nil {
			b.Fatal(err)
		}
		if newVal != "bar_"+strI {
			b.Fatal("Set/get not equal")
		}
	}
}

func Benchmark_redisTTL(b *testing.B) {
	redisCli, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err = redisCli.Cmd("SET", "key_"+strI, "bar_"+strI, "EX", 2).Err
		if err != nil {
			b.Fatal(err)
		}

		newVal, err := redisCli.Cmd("GET", "key_"+strI).Str()
		if err != nil {
			b.Fatal(err)
		}
		if newVal != "bar_"+strI {
			b.Fatal("Set/get not equal")
		}
	}
}

func Benchmark_memcacheTTL(b *testing.B) {
	mc := memcache.New("localhost:11211")

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := mc.Set(&memcache.Item{Key: "key_" + strI, Value: []byte("bar_" + strI), Expiration: 2})
		if err != nil {
			b.Fatal(err)
		}

		item, err := mc.Get("key_" + strI)
		if err != nil {
			b.Fatal(err)
		}
		newVal := item.Value
		if string(newVal) != "bar_"+strI {
			b.Fatal("Set/get not equal")
		}
	}
}
