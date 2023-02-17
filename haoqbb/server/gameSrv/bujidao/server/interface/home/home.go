package home

type IHome interface {
}

var agent IHome

func SetAgent(h IHome) {
	agent = h
}
