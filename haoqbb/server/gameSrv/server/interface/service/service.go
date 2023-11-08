package service

type IService interface {
	GetName() string // 无需重写
	GetId() int      // 无需重写
	GetLoginSrvId() int
}

type service struct {
	IService
}

var s = service{}

func SetServiceAgent(is IService) {
	s.IService = is
}

func GetServiceName() string {
	return s.GetName()
}

func GetServiceId() int {
	return s.GetId()
}

func GetLoginSrvId() int {
	return s.GetLoginSrvId()
}
