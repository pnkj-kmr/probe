package M

// LoginProfile help to describe the connection detail about
// SNMP/ SSH / HTTP etc.

type LinkedProfile struct {
	Protocol       ProtocolType `json:"protocol_type" validate:"required"`
	LoginProfileId string       `json:"login_profileid" validate:"required"`
}

type AuthSNMP struct {
	Name            string        `json:"name" validate:"required"`
	Protocol        ProtocolType  `json:"protocol_type" validate:"required,oneof=snmp"`
	Version         SNMPVersion   `json:"version" validate:"required"`
	LoginProfileId  string        `json:"login_profileid" validate:"required"`
	Community       string        `json:"community,omitempty"`
	Port            int           `json:"port,omitempty"`
	Retries         int           `json:"retries,omitempty"`
	Timeout         int           `json:"timeout,omitempty"`
	WCommunity      string        `json:"wcommunity,omitempty"`
	MaxOids         int           `json:"max_oids,omitempty"`
	Method          SNMPOperation `json:"request_type,omitempty"`
	SecurityLevel   string        `json:"security_level,omitempty"`
	AuthType        string        `json:"auth_type,omitempty"`
	AuthPass        string        `json:"auth_pass,omitempty"`
	PrivType        string        `json:"priv_type,omitempty"`
	PrivPass        string        `json:"priv_pass,omitempty"`
	Username        string        `json:"username,omitempty"`
	ContextName     string        `json:"context_name,omitempty"`
	ContextEngineID string        `json:"context_engine_id,omitempty"`
}

type AuthHTTP struct {
	Name           string       `json:"name" validate:"required"`
	Protocol       ProtocolType `json:"protocol_type" validate:"required,oneof=http https"`
	LoginProfileId string       `json:"login_profileid" validate:"required"`
	Url            string       `json:"url,omitempty"`
	Type           string       `json:"type,omitempty"`
	Host           string       `json:"host,omitempty"`
	Port           string       `json:"port,omitempty"`
	Retries        int          `json:"retries,omitempty"`
	Timeout        int          `json:"timeout,omitempty"`
	Username       string       `json:"username,omitempty"`
	Password       string       `json:"password,omitempty"`
	NoOfCustomPair int          `json:"no_of_custom_pairs,omitempty"`
	CredentialType string       `json:"credentialtype,omitempty"`
}

type AuthSFTP struct {
	Name           string       `json:"name" validate:"required"`
	Protocol       ProtocolType `json:"protocol_type" validate:"required,oneof=sftp"`
	LoginProfileId string       `json:"login_profileid" validate:"required"`
	Port           string       `json:"port,omitempty"`
	Retries        int          `json:"retries,omitempty"`
	Timeout        int          `json:"timeout,omitempty"`
	CredentialType string       `json:"credentialtype,omitempty"`
	Username       string       `json:"username,omitempty"`
	Password       string       `json:"password,omitempty"`
}

type AuthSSH struct {
	Name           string          `json:"name" validate:"required"`
	Protocol       ProtocolType    `json:"protocol_type" validate:"required,oneof=ssh"`
	LoginProfileId string          `json:"login_profileid" validate:"required"`
	Port           string          `json:"port,omitempty"`
	Retries        int             `json:"retries,omitempty"`
	Timeout        int             `json:"timeout,omitempty"`
	CredentialType string          `json:"credentialtype,omitempty"`
	Username       string          `json:"username,omitempty"`
	Password       string          `json:"password,omitempty"`
	LinkedProfiles []LinkedProfile `json:"linked_profiles,omitempty"`
}

type AuthTELNET struct {
	Name           string          `json:"name" validate:"required"`
	Protocol       ProtocolType    `json:"protocol_type" validate:"required,oneof=telnet"`
	LoginProfileId string          `json:"login_profileid" validate:"required"`
	Port           string          `json:"port,omitempty"`
	Retries        int             `json:"retries,omitempty"`
	Timeout        int             `json:"timeout,omitempty"`
	CredentialType string          `json:"credentialtype,omitempty"`
	Username       string          `json:"username,omitempty"`
	Password       string          `json:"password,omitempty"`
	LinkedProfiles []LinkedProfile `json:"linked_profiles,omitempty"`
}
