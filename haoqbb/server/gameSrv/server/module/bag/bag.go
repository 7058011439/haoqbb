package bag

import "github.com/7058011439/haoqbb/Log"

type Bag struct {
	UserId   int
	ItemList map[int]int
}

func (b *Bag) DataOK() bool {
	return b.UserId != 0
}

func (b *Bag) OnLoadEnd() {

}

func (b *Bag) Condition() map[string]interface{} {
	return map[string]interface{}{
		"userid": b.UserId,
	}
}

func (b *Bag) GiveItem(itemId int, itemCount int) {
	b.ItemList[itemId] += itemCount
	b.Update()
}

func (b *Bag) TakeItem(itemId int, itemCount int) {
	if count, ok := b.ItemList[itemId]; ok {
		if count < itemCount {
			b.ItemList[itemId] = 0
			Log.Error("减少物品失败, 物品数量不足, userId = %v, itemId = %v, 现有 = %v, 扣除 = %v", b.UserId, itemId, count, itemCount)
		} else {
			b.ItemList[itemId] += itemCount
		}
		b.Update()
	} else {
		Log.Error("减少物品失败, 物品数量不存在, userId = %v, itemId = %v, 现有 = %v, 扣除 = %v", b.UserId, itemId, count, itemCount)
	}
}

func (b *Bag) CheckItem(itemId int, itemCount int) bool {
	return b.ItemList[itemId] >= itemCount
}

func (b *Bag) Update() {
	agent.UpdateData(b.UserId, "", nil)
}
