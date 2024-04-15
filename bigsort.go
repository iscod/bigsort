package bigsort

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type BigSorted struct {
	client *redis.Client
	Key    string
	RList  []*RMember    //sort rank list
	RMap   map[any]int64 // member => rank
}

func New(key string) *BigSorted {
	s := &BigSorted{
		//client:  client,
		Key:   key,
		RList: make([]*RMember, 0),
		RMap:  make(map[any]int64, 0),
	}
	return s
}

func (s *BigSorted) ZAdd(member any, value string) error {
	score, err := decimal.NewFromString(value)
	if err != nil {
		return err
	}

	if i, ok := s.RMap[member]; ok {
		s.RList[i].Score = score
	} else {
		s.RList = append(s.RList, &RMember{Member: member, Score: score})
		s.RMap[member] = int64(len(s.RList) - 1)
	}

	s.sort(member) //插入排序
	return nil
}

// list fast for v1
func (s *BigSorted) ZIncrBy(member any, value string) error {
	score, err := decimal.NewFromString(value)
	if err != nil {
		return err
	}
	//s.client.HSet(context.TODO(), member, value)

	if i, ok := s.RMap[member]; ok {
		s.RList[i].Score = s.RList[i].Score.Add(score)
	} else {
		s.RList = append(s.RList, &RMember{Member: member, Score: score})
		s.RMap[member] = int64(len(s.RList) - 1)
	}

	s.sort(member) //插入排序
	return nil
}

func (s *BigSorted) Count() int {
	return len(s.RMap)
}

func (s *BigSorted) sort(member any) {
	r, ok := s.RMap[member]
	if !ok {
		return
	}
	score := s.RList[r].Score

	for i := int(r); i > 0; i-- { //向前移动,两者只会发生一个
		if s.RList[i-1].Score.LessThan(score) {
			s.RList[i], s.RList[i-1] = s.RList[i-1], s.RList[i]
		} else {
			break
		}
	}

	for i := int(r); i < len(s.RList)-1; i++ { //向后移动
		if s.RList[i+1].Score.GreaterThan(score) {
			s.RList[i], s.RList[i+1] = s.RList[i+1], s.RList[i]
		} else {
			break
		}
	}
}

// todo
func (s *BigSorted) Remove(member string) {
	if i, ok := s.RMap[member]; ok {
		delete(s.RMap, member)
		s.RList = append(s.RList[:i], s.RList[i+1:]...)
	}
}

func (s *BigSorted) ZRevRank(start, stop int) []*RMember {
	if stop < 0 {
		return s.RList[start:stop]
	}
	if stop > len(s.RList) {
		return s.RList[start:]
	} else {
		return s.RList[start:stop]
	}
}

func (s *BigSorted) ZRank(member string) int64 {
	if r, ok := s.RMap[member]; ok {
		return r
	} else {
		return -1
	}
}

func (s *BigSorted) ZRem(member any) *RMember {
	if r, ok := s.RMap[member]; ok {
		return s.RList[r]
	} else {
		return nil
	}
}

// From redis data Unmarshal BigSorted
func (s *BigSorted) From(client *redis.Client) {
	con := context.TODO()
	iter := client.ZScan(con, s.Key, 0, "", 0).Iterator()
	//iter := client.HScan(con, key, 0, "", 1000).Iterator()
	//iter := client.HScan(con, key, 0, "", 1000).Iterator()
	for iter.Next(con) {
		field := iter.Val()
		iter.Next(con)
		value := iter.Val()
		s.ZAdd(field, value)
	}
}
