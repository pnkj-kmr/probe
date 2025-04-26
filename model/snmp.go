package M

type SNMPVersion string
type SNMPOperation uint8

const (
	GET      SNMPOperation = 0x1
	WALK     SNMPOperation = 0x2
	BULKWALK SNMPOperation = 0x3

	VERSION1  SNMPVersion = "1"
	VERSION2C SNMPVersion = "2c"
	VERSION3  SNMPVersion = "3"
)

type SNMPParams struct {
	PollPeriod int      `json:"poll_period"`
	PollType   string   `json:"poll_type"`
	PCid       string   `json:"parent_ci_id"`
	MibProfile string   `json:"mib_profile"`
	IP         string   `json:"poll_addr"`
	IPs        []string `json:"poll_addr_list,omitempty"`

	CiName     string `json:"ci_name"`
	CiCategory string `json:"ci_category"`
	Speed      int    `json:"speed"`
	Memory     int    `json:"memory_size"`
	// OidIndex   string `json:"oidIndex"`
	CustomType string `json:"custom_type"` //"MAC"|"STRING"|Hex-STRING
}

type SNMPReq struct {
	IP             string      `json:"poll_addr"`
	Cid            string      `json:"ci_id"`
	Params         SNMPParams  `json:"params,omitempty"`
	InputStats     []InputStat `json:"input_stats,omitempty"`
	LoginProfileID string      `json:"login_profileid"`
}

type SNMPRes struct {
	Cid    string       `json:"ci_id"`
	Params SNMPParams   `json:"params,omitempty"`
	T      string       `json:"t,omitempty"`
	Stats  []OutputStat `json:"stats"`
	Err    string       `json:"error"`
}
