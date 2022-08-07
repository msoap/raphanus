/*
run:

	docker run --name raphanus --rm --publish 8771:8771 msoap/raphanus
	docker run --name redis --rm --publish 6379:6379 redis
	docker run --name memcache --rm --publish 11211:11211 memcached
	make run-benchmark
*/
package raphanus_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/msoap/raphanus"
	raphanusclient "github.com/msoap/raphanus/client"
	"github.com/stretchr/testify/require"
)

func Benchmark_raphanusServer(b *testing.B) {
	raph := raphanusclient.New(raphanusclient.Cfg{Address: "http://localhost:8771"})

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := raph.SetStr("key_"+strI, "bar_"+strI, 0)
		require.NoError(b, err)

		newVal, err := raph.GetStr("key_" + strI)
		require.NoError(b, err)
		require.Equal(b, newVal, "bar_"+strI)
	}
}

func Benchmark_raphanusEmbed(b *testing.B) {
	raph := raphanus.New()

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := raph.SetStr("key_"+strI, "bar_"+strI, 0)
		require.NoError(b, err)

		newVal, err := raph.GetStr("key_" + strI)
		require.NoError(b, err)
		require.Equal(b, newVal, "bar_"+strI)
	}
}

func Benchmark_boltdb(b *testing.B) {
	db, err := bolt.Open("bolt_bench_tmp.db", 0600, nil)
	require.NoError(b, err)
	defer func() {
		require.NoError(b, db.Close())
		require.NoError(b, os.Remove("bolt_bench_tmp.db"))
	}()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("bucket"))
		return err
	})
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)

		err := db.Update(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte("bucket")).Put([]byte("key_"+strI), []byte("bar_"+strI))
		})
		require.NoError(b, err)

		err = db.View(func(tx *bolt.Tx) error {
			newVal := tx.Bucket([]byte("bucket")).Get([]byte("key_" + strI))
			require.Equal(b, string(newVal), "bar_"+strI)
			return nil
		})
		require.NoError(b, err)
	}
}

func Benchmark_redis(b *testing.B) {
	redisCli, err := redis.Dial("tcp", "localhost:6379")
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err = redisCli.Cmd("SET", "key_"+strI, "bar_"+strI).Err
		require.NoError(b, err)

		newVal, err := redisCli.Cmd("GET", "key_"+strI).Str()
		require.NoError(b, err)
		require.Equal(b, newVal, "bar_"+strI)
	}
}

func Benchmark_memcache(b *testing.B) {
	mc := memcache.New("localhost:11211")

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := mc.Set(&memcache.Item{Key: "key_" + strI, Value: []byte("bar_" + strI)})
		require.NoError(b, err)

		item, err := mc.Get("key_" + strI)
		require.NoError(b, err)
		newVal := item.Value
		require.Equal(b, string(newVal), "bar_"+strI)
	}
}

func Benchmark_raphanusServerTTL(b *testing.B) {
	raph := raphanusclient.New(raphanusclient.Cfg{Address: "http://localhost:8771"})

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := raph.SetStr("key_"+strI, "bar_"+strI, 2)
		require.NoError(b, err)

		newVal, err := raph.GetStr("key_" + strI)
		if err != nil {
			// skip deleted keys error
			continue
		}
		require.Equal(b, newVal, "bar_"+strI)
	}
}

func Benchmark_raphanusEmbedTTL(b *testing.B) {
	raph := raphanus.New()

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := raph.SetStr("key_"+strI, "bar_"+strI, 2)
		require.NoError(b, err)

		newVal, err := raph.GetStr("key_" + strI)
		require.NoError(b, err)
		require.Equal(b, newVal, "bar_"+strI)
	}
}

func Benchmark_redisTTL(b *testing.B) {
	redisCli, err := redis.Dial("tcp", "localhost:6379")
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err = redisCli.Cmd("SET", "key_"+strI, "bar_"+strI, "EX", 2).Err
		require.NoError(b, err)

		newVal, err := redisCli.Cmd("GET", "key_"+strI).Str()
		require.NoError(b, err)
		require.Equal(b, newVal, "bar_"+strI)
	}
}

func Benchmark_memcacheTTL(b *testing.B) {
	mc := memcache.New("localhost:11211")

	for i := 0; i < b.N; i++ {
		strI := strconv.Itoa(i)
		err := mc.Set(&memcache.Item{Key: "key_" + strI, Value: []byte("bar_" + strI), Expiration: 2})
		require.NoError(b, err)

		item, err := mc.Get("key_" + strI)
		require.NoError(b, err)
		newVal := item.Value
		require.Equal(b, string(newVal), "bar_"+strI)
	}
}
