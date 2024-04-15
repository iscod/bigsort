package bigsort

import (
	"github.com/shopspring/decimal"
	"math"
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkZadd(b *testing.B) {
	// run the Zadd function b.N times
	maxInt64 := decimal.NewFromInt(math.MaxInt64)
	bigSort := New("redis-hset-key")
	for i := 0; i < 100000; i++ {
		bigSort.ZAdd(i, maxInt64.Add(decimal.NewFromInt(int64(i))).String())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bigSort.ZAdd(n, maxInt64.Add(decimal.NewFromInt(int64(100000+rand.Intn(b.N)))).String())
	}
}

func BenchmarkZRem(b *testing.B) {
	// run the ZRem function b.N times
	maxInt64 := decimal.NewFromInt(math.MaxInt64)
	bigSort := New("redis-hset-key")
	for i := 0; i < 100000; i++ {
		bigSort.ZAdd(i, maxInt64.Add(decimal.NewFromInt(int64(i))).String())
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bigSort.ZRem(n)
	}
}

func TestZRem(t *testing.T) {
	bigSort := New("redis-hset-key")
	//_ = bigSort.ZAdd("member1", "922337203685477580800000001") // > math.MaxInt64
	//_ = bigSort.ZAdd("member2", "922337203685477580800000002") // > math.MaxInt64

	score := decimal.NewFromInt(math.MaxInt64).Add(decimal.NewFromInt(int64(rand.Intn(math.MaxInt))))
	_ = bigSort.ZAdd("member3", score.String()) // > math.MaxInt64

	member := bigSort.ZRem("member3") //rank == 2
	if !member.Score.Equal(score) {
		t.Errorf("Test zrem error socre %s, zrem value %s", score.String(), member.Score.String())
	}
}

func TestRank(t *testing.T) {
	bigSort := New("redis-hset-key")
	_ = bigSort.ZAdd("member1", "922337203685477580800000001") // > math.MaxInt64
	_ = bigSort.ZAdd("member2", "922337203685477580800000002") // > math.MaxInt64
	_ = bigSort.ZAdd("member3", "922337203685477580800000003") // > math.MaxInt64

	rank := bigSort.ZRank("member3") //rank == 2
	if rank != 2 {
		t.Errorf("Test rank = %d, want %d", rank, 2)
	}
}

func TestFromRedis(t *testing.T) {
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr:     "127.0.0.1:6379",
	//	PoolSize: 200,
	//})

	memberLen := 10000

	b := New("testbignumber")
	//maxInt64 := decimal.NewFromInt(math.MaxInt64).Add(decimal.NewFromInt(math.MaxInt64))
	maxInt64 := decimal.NewFromInt(int64(memberLen))
	for i := 0; i < memberLen; i++ {
		b.ZAdd(strconv.Itoa(i), maxInt64.Add(decimal.NewFromInt(int64(i))).String())
	}
	if len(b.RList) != memberLen {
		t.Errorf("list number members len error, sortLen : %d, member : %d", len(b.RMap), memberLen)
	}

	m := rand.Intn(memberLen)
	r := b.ZRank(strconv.Itoa(m))
	if r != int64(m) {
		t.Errorf("Form Rdis rank error, rank: %d, m: %d, score : %s", r, m, b.ZRem(strconv.Itoa(m)).Score.String())
	}

}
