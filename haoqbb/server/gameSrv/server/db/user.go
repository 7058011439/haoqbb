package db

type mission struct {
	Id         int `bson:"id,omitempty"`
	CurrentPro int
	TotalPro   int
}

type missions struct {
	MissionList []*mission `bson:"mission_list,omitempty"`
	IntData     int
	StringData  string
}

type User struct {
	Id       int64     `bson:"userid,omitempty"`
	NickName string    `bson:"nickname,omitempty"`
	Mission  *missions `bson:"mission,omitempty"`
	Friends  []uint64  `bson:"friends,omitempty"`
}

func (u *User) Update() {
	agent.UpdateData(u.Id)
}

func (u *User) GetUserId() int64 {
	return u.Id
}

func (u *User) IsValid() bool {
	return u.Id != 0 && u.NickName != ""
}

func (u *User) Condition() map[string]interface{} {
	return map[string]interface{}{
		"userid": u.Id,
	}
}
