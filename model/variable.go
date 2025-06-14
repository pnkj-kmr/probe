package M

type (
	ProbeType    string
	ProtocolType string
)

// declaring generic constraints

const (
	// restful apiserver related flag
	API ProbeType = "api"
	// export related variable
	KAFKA   ProbeType = "kafka"
	WEBHOOK ProbeType = "webhook"
	DBDUMP  ProbeType = "dbdump"
	RMQ     ProbeType = "rmq"
	// db related variables
	EVENT   ProbeType = "events"
	PROCESS ProbeType = "processes"
	CONFIG  ProbeType = "configuration"
	CRED    ProbeType = "credential"
	ENV     ProbeType = "environment"
	// db and polling related variables
	ICMP   ProbeType = "icmp"
	SNMP   ProbeType = "snmp"
	HTTP   ProbeType = "http"
	SSH    ProbeType = "ssh"
	TELNET ProbeType = "telnet"
)

// interval related constants
// db, poll, schedule etc
const (
	INTERVAL_60  int = 60
	INTERVAL_300 int = 300
)

const (
	SNMPType   ProtocolType = "snmp"
	HTTPType   ProtocolType = "http"
	HTTPSType  ProtocolType = "https"
	SSHType    ProtocolType = "ssh"
	SFTPType   ProtocolType = "sftp"
	TELNETType ProtocolType = "telnet"
)
