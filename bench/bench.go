package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Simple benchmark test. Generate random keys/vals and insert into the DB. For some % of these, wait for some
// amount of time and then query for them.

const (
	minChar = ' '
	maxChar = '~'
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randStr(n int) string {
	r := make([]rune, n)
	for i := range r {
		r[i] = rune(rand.Intn(maxChar-minChar+1) + minChar)
	}
	return string(r)
}

func makeRequests(pool *redis.Pool) {
	conn := pool.Get()
	defer conn.Close()

	const N = 10000

	start := time.Now()
	for i := 0; i < N; i++ {
		key := randStr(10)
		val := randStr(1000)
		_, err := conn.Do("SET", key, val)
		if err != nil {
			log.Fatal(err)
		}
		//result, err := redis.String(conn.Do("GET", key))
		//if err != nil {
		//log.Fatal(err)
		//}
		//if result != val {
		//log.Fatal("result mismatch")
		//}
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("Took %s for %d iterations | %s / op | %.1f ops / s\n", elapsed, N, elapsed/N,
		N*float64(time.Second)/float64(elapsed))
}

func main() {
	pool := &redis.Pool{
		MaxIdle:     0,
		MaxActive:   0, // No limit
		IdleTimeout: time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "localhost:5533") },
	}
	makeRequests(pool)
}
