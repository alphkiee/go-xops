package request

// 适用于前端传过来的
type IdsReq struct {
	Ids []uint `json:"ids" form:"ids"` // 传多个id
}

// 适用于前端传过来的
type KeyReq struct {
	Key string `json:"key" form:"key"` // 传多个id
}
