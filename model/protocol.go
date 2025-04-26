package M

// LoginProfile help to describe the connection detail about
// SNMP/ SSH / HTTP etc.

type ProtocolType string

const (
	SNMP   ProtocolType = "snmp"
	HTTP   ProtocolType = "http"
	SSH    ProtocolType = "ssh"
	SFTP   ProtocolType = "sftp"
	TELNET ProtocolType = "telnet"
)

type LinkedProfile struct {
	Protocol       ProtocolType `json:"protocol_type"`
	LoginProfileId string       `json:"login_profileid"`
}

type LoginProfile struct {
	Name           string       `json:"name"`
	Protocol       ProtocolType `json:"protocol_type"`
	LoginProfileId string       `json:"login_profileid"`
	Port           int          `json:"port,omitempty"`
	Retries        int          `json:"retries,omitempty"`
	Timeout        int          `json:"timeout,omitempty"`
	Username       string       `json:"username,omitempty"`
	// snmp v2c related params
	Version    SNMPVersion   `json:"version,omitempty"`
	Community  string        `json:"community,omitempty"`
	WCommunity string        `json:"wcommunity,omitempty"`
	Method     SNMPOperation `json:"request_type,omitempty"` // get back on this
	// snmp v3 related params
	SecurityLevel   string `json:"securitylevel,omitempty"`
	AuthType        string `json:"authtype,omitempty"`
	AuthPass        string `json:"authpass,omitempty"`
	PrivType        string `json:"privtype,omitempty"`
	PrivPass        string `json:"privpass,omitempty"`
	ContextName     string `json:"context_name,omitempty"`
	ContextEngineID string `json:"context_engine_id,omitempty"`
	// telnet/sftp/ssh/http params
	Password       string `json:"password,omitempty"`
	NoOfCustomPair int    `json:"no_of_custom_pairs,omitempty"`
	CredentialType string `json:"credentialtype,omitempty"`
	// ssh params
	LinkedProfiles []LinkedProfile `json:"linked_profiles,omitempty"`
	// http params
	Type string `json:"type,omitempty"`
}
