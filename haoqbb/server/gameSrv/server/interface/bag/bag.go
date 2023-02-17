package bag

type IBag interface {
	GiveItem(userId int, itemId int, itemCount int)
	TakeItem(userId int, itemId int, itemCount int)
	CheckItem(userId int, itemId int, itemCount int) bool
}

var agent IBag

func SetAgent(b IBag) {
	agent = b
}

func GiveItem(userId, itemId, itemCount int) {
	agent.GiveItem(userId, itemId, itemCount)
}

func TakeItem(userId, itemId, itemCount int) {
	agent.TakeItem(userId, itemId, itemCount)
}

func CheckItem(userId, itemId, itemCount int) bool {
	return agent.CheckItem(userId, itemId, itemCount)
}
