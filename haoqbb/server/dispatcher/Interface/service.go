package Interface

type IService interface {
	GetName() string // 无需重写
	GetId() int      // 无需重写
}

type service struct {
	IService
}

var s = service{}

func SetServiceAgent(is IService) {
	s.IService = is
}

func GetName() string {
	return s.GetName()
}

func GetId() int {
	return s.GetId()
}
