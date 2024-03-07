package player

type Player struct {
	UserId int
	OpenId string
}

func (p *Player) DataOK() bool {
	return p.UserId != 0
}

func (p *Player) OnLoadEnd() {

}

func (p *Player) Condition() map[string]interface{} {
	return map[string]interface{}{
		"userid": p.UserId,
	}
}
