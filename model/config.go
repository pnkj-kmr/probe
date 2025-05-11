package M

// declaring generic constraints
const (
	_ = iota
	API
	KAFKA
	WEBHOOK
	RMQ
	CONFIG
	AUTH_PROFILE
	ENV
	ICMP
	SNMP
	HTTP
	SSH
	TELNET
)

// interval related constants
// db, poll, schedule etc
const (
	INTERVAL_60  int = 5  //60		// added for testing as 5 seconds
	INTERVAL_300 int = 15 //300
)

type ProtocolType string

const (
	SNMPType   ProtocolType = "snmp"
	HTTPType   ProtocolType = "http"
	HTTPSType  ProtocolType = "https"
	SSHType    ProtocolType = "ssh"
	SFTPType   ProtocolType = "sftp"
	TELNETType ProtocolType = "telnet"
)
