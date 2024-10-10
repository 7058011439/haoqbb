package db

import "github.com/7058011439/haoqbb/Log"

type Bag struct {
	IDBData  `bson:"-"`
	ItemList map[int]int
}

func NewBag() *Bag {
	return &Bag{
		ItemList: map[int]int{},
	}
}

func (b *Bag) Copy() *Bag {
	ret := &Bag{
		ItemList: map[int]int{},
	}
	for k, v := range b.ItemList {
		ret.ItemList[k] = v
	}
	return ret
}

func (b *Bag) GiveItem(itemId int, itemCount int) {
	b.ItemList[itemId] += itemCount
	b.Update()
}

func (b *Bag) TakeItem(itemId int, itemCount int) {
	if count, ok := b.ItemList[itemId]; ok {
		if count < itemCount {
			b.ItemList[itemId] = 0
			Log.Error("减少物品失败, 物品数量不足, userId = %v, itemId = %v, 现有 = %v, 扣除 = %v", b.GetUserId(), itemId, count, itemCount)
		} else {
			b.ItemList[itemId] -= itemCount
		}
		b.Update()
	} else {
		Log.Error("减少物品失败, 物品id不存在, userId = %v, itemId = %v, 现有 = %v, 扣除 = %v", b.GetUserId(), itemId, count, itemCount)
	}
}

func (b *Bag) CheckItem(itemId int, itemCount int) bool {
	if count, ok := b.ItemList[itemId]; ok {
		return count >= itemCount
	} else {
		return false
	}
}
