package bag

type IBag interface {
	UseItem(userId int, itemId int, itemCount int) bool
}

var agent IBag

func SetAgent(b IBag) {
	agent = b
}

func UseItem(userId int, itemId int, itemCount int) bool {
	return agent.UseItem(userId, itemId, itemCount)
}
