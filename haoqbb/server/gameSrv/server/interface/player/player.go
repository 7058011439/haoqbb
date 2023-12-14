package player

type IPlayer interface {
	Kick(userId int64)
	GetClientId(userId int64) uint64
	GetUserId(clientId uint64) int64
	Login(clientId uint64, userId int64)
	LogOut(clientId uint64, userId int64)
}

type SPlayer struct {
}

func (s *SPlayer) Kick(userId int64) {

}

var agent IPlayer

func SetAgent(b IPlayer) {
	agent = b
}

func Kick(userId int64) {
	agent.Kick(userId)
}

func GetClientId(userId int64) uint64 {
	return agent.GetClientId(userId)
}

func GetUserId(clientId uint64) int64 {
	return agent.GetUserId(clientId)
}

func Login(clientId uint64, userId int64) {
	agent.Login(clientId, userId)
}

func LogOut(clientId uint64, userId int64) {
	agent.LogOut(clientId, userId)
}
