package home

type IPlayer interface {
}

var agent IPlayer

func SetAgent(p IPlayer) {
	agent = p
}
