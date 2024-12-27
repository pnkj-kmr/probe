package model

type PingReq struct {
	IP  string `json:"ip"`
	Cid string `json:"ci_id"`
}

type PingRes struct {
	IP  string `json:"ip"`
	Cid string `json:"ci_id"`
}
