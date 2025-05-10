package snmp

import (
	"encoding/hex"
	"fmt"
	"log"
	M "probe/model"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	g "github.com/gosnmp/gosnmp"
)

func doSNMP(i M.SNMPReq, p M.AuthSNMP) (out M.SNMPRes, err error) {
	out.Cid = i.Cid
	out.Params = i.Params

	_snmp := getSNMPInstance(p)
	err = _snmp.Connect()
	if err != nil {
		out.Err = err.Error()
		log.Println("SNMP connection err", err)
		return
	}
	defer _snmp.Conn.Close()

	var data []M.OutputStat
	switch p.Method {
	case M.BULKWALK:
		data, err = snmpBulkWalk(_snmp, i.InputStats, i.Params.CustomType)
		if err != nil {
			out.Err = err.Error()
			return
		}
	case M.WALK:
		data, err = snmpWalk(_snmp, i.InputStats, i.Params.CustomType)
		if err != nil {
			out.Err = err.Error()
			return
		}
	default:
		data, err = snmpGet(_snmp, i.InputStats, i.Params.CustomType)
		if err != nil {
			out.Err = err.Error()
			return
		}
	}

	out.Stats = data
	out.T = time.Now().UTC().Format("2006-01-02 15:04:05.000000")
	return
}

func getSNMPInstance(p M.AuthSNMP) *g.GoSNMP {
	//
	//	snmpwalk -v 1 <target> <oid (i.e. 1.3.6.1.2.1.1.3.0)>
	//	snmpwalk -v 2c -c <community> <target> <oid (i.e. 1.3.6.1.2.1.1.3.0)>
	// 	snmpwalk -v 3 -l <level> -u <username> -a <authtype> -x <privtype> -A <authpass> -X <privpass>  <target> <oid>
	//

	var _port uint16 = uint16(p.Port)
	if _port == 0 {
		_port = 161
	}
	var community string = p.Community
	if community == "" {
		community = "public"
	}
	var timeout int = p.Timeout
	if timeout == 0 {
		timeout = 5
	}
	var msgFlag string = p.SecurityLevel
	if msgFlag == "" {
		msgFlag = "AuthPriv"
	}
	// we are not allowing more than 3 retries of a IP
	// helps to reducs the poll cycle
	var retries = p.Retries
	if retries > 3 {
		retries = 3
	}
	var maxOids = 60
	if p.MaxOids > maxOids {
		maxOids = p.MaxOids
	}

	version := getVersion(p.Version)

	_snmp := g.Default

	_snmp.Target = ""
	_snmp.Port = _port
	_snmp.Retries = retries
	_snmp.Timeout = time.Duration(timeout) * time.Second
	_snmp.MaxOids = maxOids
	_snmp.ExponentialTimeout = true
	_snmp.Version = version

	switch version {
	case g.Version1:
		_snmp.Transport = "udp"
		_snmp.Community = p.Community
	case g.Version2c:
		_snmp.Transport = "udp"
		_snmp.Community = p.Community
	case g.Version3:
		_snmp.SecurityModel = g.UserSecurityModel
		_snmp.MsgFlags = getMsgFlag(msgFlag)
		_snmp.SecurityParameters = &g.UsmSecurityParameters{
			UserName:                 p.Username,
			AuthenticationProtocol:   getAuthType(p.AuthType),
			AuthenticationPassphrase: p.AuthPass,
			PrivacyProtocol:          getPrivType(p.PrivType),
			PrivacyPassphrase:        p.PrivPass,
		}
		_snmp.ContextName = p.ContextName
		if p.ContextEngineID != "" {
			_snmp.ContextEngineID = p.ContextEngineID
		}
	}

	return _snmp
}

func snmpGet(_snmp *g.GoSNMP, stats []M.InputStat, customType string) (data []M.OutputStat, err error) {
	var oids []string
	var oidMap = make(map[string]string)
	for _, o := range stats {
		oidMap[o.Oid] = o.Dn
		oids = append(oids, o.Oid)
	}

	result, err := _snmp.Get(oids)
	if err != nil {
		log.Println("SNMP GET err", err)
		return
	}
	for _, v := range result.Variables {
		data = append(data, parseData(v, oidMap, customType))
	}
	return
}

func snmpWalk(_snmp *g.GoSNMP, stats []M.InputStat, customType string) (data []M.OutputStat, err error) {
	// size calculation
	//for one value in []M.OutputStat 16 bytes(interface{}) + 16 bytes (string) + 16 bytes (string) = 48 bytes
	// for 1000 it will be 48*100 = 4800 bytes and plus slice header space is 24 bytes
	// so overall ~4.6Kb
	var oids []string
	var oidMap = make(map[string]string)
	for _, o := range stats {
		oidMap[o.Oid] = o.Dn
		oids = append(oids, o.Oid)
	}

	data = make([]M.OutputStat, 0, 300) // preallocation of memory help in performance boost. As memory reallocation happens less
	var callback = func(d g.SnmpPDU) error {
		data = append(data, parseData(d, oidMap, customType))
		return nil
	}
	err = _snmp.Walk(oids[0], callback)
	if err != nil {
		log.Println("SNMP WALK err", err)
		return
	}
	return
}

func snmpBulkWalk(_snmp *g.GoSNMP, stats []M.InputStat, customType string) (data []M.OutputStat, err error) {
	var oids []string
	var oidMap = make(map[string]string)
	for _, o := range stats {
		oidMap[o.Oid] = o.Dn
		oids = append(oids, o.Oid)
	}

	data = make([]M.OutputStat, 0, 300) // preallocation of memory help in performance boost. As memory reallocation happens less
	var callback = func(d g.SnmpPDU) error {
		data = append(data, parseData(d, oidMap, customType))
		return nil
	}
	err = _snmp.BulkWalk(oids[0], callback)
	if err != nil {
		log.Println("SNMP BULK WALK err", err)
		return
	}
	return
}

func parseData(d g.SnmpPDU, oidMap map[string]string, customType string) (data M.OutputStat) {
	dn := oidMap[d.Name]
	switch d.Type {
	case g.Integer, g.Counter32, g.Gauge32, g.TimeTicks, g.Counter64, g.Uinteger32:
		data = M.OutputStat{Dn: dn, Oid: d.Name, Value: g.ToBigInt(d.Value), Type: pduTypeToString(d.Type)}
	case g.OctetString:
		b, ok := d.Value.([]byte)
		if !ok || b == nil {
			data = M.OutputStat{Dn: dn, Oid: d.Name, Value: "", Type: "STRING"}
			break
		}
		value, new_type := parseOctetString(b, customType)
		data = M.OutputStat{Dn: dn, Oid: d.Name, Value: value, Type: new_type}
	default:
		// default value should be string and type should be specified
		strVal, ok := d.Value.(string)
		if !ok {
			strVal = fmt.Sprintf("%v", d.Value)
		}
		data = M.OutputStat{Dn: dn, Oid: d.Name, Value: strVal, Type: pduTypeToString(d.Type)}
	}
	return
}

func isMostlyPrintable(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	printableCount := 0
	for _, b := range data {
		if unicode.IsPrint(rune(b)) {
			printableCount++
		}
	}
	return float64(printableCount)/float64(len(data)) > 0.8
}

func formHexBytesWithSpaces(data []byte) string {
	parts := make([]string, len(data))
	for i, b := range data {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, " ")
}

func parseOctetString(value []byte, customType string) (string, string) {

	if customType == "MAC" {
		var macAddress strings.Builder
		for i, b := range value {
			if i > 0 {
				macAddress.WriteString(":")
			}
			macAddress.WriteString(fmt.Sprintf("%02X", b))
		}
		return macAddress.String(), "MAC"
	} else if customType == "STRING" {
		return string(value), "STRING"
	}

	if isMostlyPrintable(value) && utf8.Valid(value) {
		str := string(value)
		if decoded, err := hex.DecodeString(str); err == nil && isMostlyPrintable(decoded) {
			return string(decoded), "Hex-STRING"
		}
		return str, "STRING"
	}

	if len(value) == 6 {
		var macAddress strings.Builder
		for i, b := range value {
			if i > 0 {
				macAddress.WriteString(":")
			}
			macAddress.WriteString(fmt.Sprintf("%02X", b))
		}
		return macAddress.String(), "MAC"
	}

	return formHexBytesWithSpaces(value), "Hex-STRING"
}
