# bigsort
Helper Big numbers Sorted Set for redis

# Ask what needs bigsort？

Redis sorted sets use a double 64-bit floating point number to represent the score，that is able to represent precisely integer numbers between -(2^53) and +(2^53) included。

So bigger numbers Lose effectiveness！！！

# Adaptive 

1. The new key is usually the smallest
2. The increase in value will not cause much change

# Benchmark on 10,0000 levels

```bash
go test -bench=.
```

```
machine Apple M1 Pro 8 cpus
goos: darwin
goarch: arm64
pkg: github.com/iscod/bigsort
BenchmarkZadd-8   	   51714	    224779 ns/op
BenchmarkZRem-8   	34261527	        37.47 ns/op
```


## Installation

```bash
go get -u github.com/iscod/bigsort
```

# Example

```Go
bsort := bigsort.New("redis-hset-key")
_ = bsort.ZAdd("member1", "922337203685477580800000001") // > math.MaxInt64
_ = bsort.ZAdd("member2", "922337203685477580800000002") // > math.MaxInt64
_ = bsort.ZAdd("member3", "922337203685477580800000003") // > math.MaxInt64

rank := bsort.ZRank("member3") //rank == 2

rankList := bsort.ZRevRank(0, 50) //rank 0-50

```