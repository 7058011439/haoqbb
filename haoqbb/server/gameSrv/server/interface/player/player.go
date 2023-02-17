package player

type IPlayer interface {
	Kick(userId int)
	GetClientId(userId int) uint64
	GetUserId(clientId uint64) int
	Login(clientId uint64, userId int)
	LogOut(clientId uint64, userId int)
}

var agent IPlayer

func SetAgent(b IPlayer) {
	agent = b
}

func Kick(userId int) {
	agent.Kick(userId)
}

func GetClientId(userId int) uint64 {
	return agent.GetClientId(userId)
}

func GetUserId(clientId uint64) int {
	return agent.GetUserId(clientId)
}

func Login(clientId uint64, userId int) {
	agent.Login(clientId, userId)
}

func LogOut(clientId uint64, userId int) {
	agent.LogOut(clientId, userId)
}
