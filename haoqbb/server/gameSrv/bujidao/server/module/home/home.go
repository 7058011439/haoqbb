package home

type furniture struct {
	ItemId int32
	Pos    int32
	Dir    int32
}

type home struct {
	UserId        int
	Level         int32
	FurnitureList []*furniture
}

func (h *home) DataOK() bool {
	return h.UserId != 0
}

func (h *home) OnLoadEnd() {

}

func (h *home) Condition() map[string]interface{} {
	return map[string]interface{}{
		"userid": h.UserId,
	}
}

func (h *home) saveTheme() {

}

func (h *home) upgradeLevel() {
	h.Level++
	h.update()
}

func (h *home) update() {
	agent.UpdateData(h.UserId, "", nil)
}
