package M

type (
	ProbeType    string
	ProtocolType string
)

// declaring generic constraints

const (
	API     ProbeType = "api"
	KAFKA   ProbeType = "kafka"
	WEBHOOK ProbeType = "webhook"
	RMQ     ProbeType = "rmq"
	CONFIG  ProbeType = "configuration"
	CRED    ProbeType = "credential"
	ENV     ProbeType = "environment"
	ICMP    ProbeType = "icmp"
	SNMP    ProbeType = "snmp"
	HTTP    ProbeType = "http"
	SSH     ProbeType = "ssh"
	TELNET  ProbeType = "telnet"
)

// interval related constants
// db, poll, schedule etc
const (
	INTERVAL_60  int = 5  //60		// added for testing as 5 seconds
	INTERVAL_300 int = 15 //300
)

const (
	SNMPType   ProtocolType = "snmp"
	HTTPType   ProtocolType = "http"
	HTTPSType  ProtocolType = "https"
	SSHType    ProtocolType = "ssh"
	SFTPType   ProtocolType = "sftp"
	TELNETType ProtocolType = "telnet"
)
