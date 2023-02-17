package Stl

import (
	"Core/Util"
	"container/list"
	"fmt"
)

func NewSet() *Set {
	return &Set{
		listData: list.New(),
		mapData:  make(map[interface{}]*list.Element),
	}
}

type Set struct {
	listData *list.List
	mapData  map[interface{}]*list.Element
}

func (s *Set) Add(data interface{}) bool {
	if s.mapData[data] == nil {
		element := s.listData.PushBack(data)
		s.mapData[data] = element
		return true
	}
	return false
}

func (s *Set) Del(data interface{}) bool {
	if s.mapData[data] != nil {
		if _, ok := s.mapData[data]; ok {
			s.listData.Remove(s.mapData[data])
			delete(s.mapData, data)
			return true
		}
	}

	return false
}

func (s *Set) Empty() bool {
	return len(s.mapData) == 0
}

func (s *Set) Len() int {
	return s.listData.Len()
}

func (s *Set) String() string {
	ret := ""
	for e := s.listData.Front(); e != nil; e = e.Next() {
		ret += fmt.Sprintf("%v,", e.Value)
	}
	if s.Len() > 0 {
		ret = ret[0 : len(ret)-1]
	}
	return ret
}

func (s *Set) Range(f func(interface{})) {
	for e := s.listData.Front(); e != nil; e = e.Next() {
		f(e.Value)
	}
}

func (s *Set) Exist(data interface{}) bool {
	for e := s.listData.Front(); e != nil; e = e.Next() {
		if Util.CompareEqual(e.Value, data) {
			return true
		}
	}
	return false
}
