package M

type InputStat struct {
	Dn  string `json:"dn"`
	Oid string `json:"oid"`
}

type OutputStat struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
	Oid   string `json:"oid"`
	Dn    string `json:"dn"`
}
