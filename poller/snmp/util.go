package snmp

import (
	"fmt"
	M "probe/model"

	g "github.com/gosnmp/gosnmp"
)

func getVersion(v M.SNMPVersion) g.SnmpVersion {
	switch v {
	case M.VERSION1:
		return g.Version1
	case M.VERSION2C:
		return g.Version2c
	case M.VERSION3:
		return g.Version3
	default:
		return g.Version2c
	}
}

func getMsgFlag(s string) g.SnmpV3MsgFlags {
	switch s {
	case "NoAuthNoPriv", "noAuthNoPriv":
		return g.NoAuthNoPriv
	case "AuthNoPriv", "authNoPriv":
		return g.AuthNoPriv
	case "AuthPriv", "authPriv":
		return g.AuthPriv
	case "Reportable":
		return g.Reportable
	default:
		return g.NoAuthNoPriv
	}
}

func getAuthType(s string) g.SnmpV3AuthProtocol {
	switch s {
	case "MD5":
		return g.MD5
	case "SHA":
		return g.SHA
	case "SHA224", "SHA-224":
		return g.SHA224
	case "SHA256", "SHA-256":
		return g.SHA256
	case "SHA384", "SHA-384":
		return g.SHA384
	case "SHA512", "SHA-512":
		return g.SHA512
	default:
		return g.NoAuth
	}
}

func getPrivType(s string) g.SnmpV3PrivProtocol {
	switch s {
	case "DES":
		return g.DES
	case "AES", "AES-128":
		return g.AES
	case "AES192", "AES-192":
		return g.AES192
	case "AES256", "AES-256":
		return g.AES256
	case "AES192C", "AES-192C":
		return g.AES192C
	case "AES256C", "AES-256C":
		return g.AES256C
	default:
		return g.NoPriv
	}
}

func pduTypeToString(t g.Asn1BER) string {
	switch t {
	case g.Boolean:
		return "Boolean"
	case g.Integer:
		return "Integer"
	case g.BitString:
		return "BitString"
	case g.OctetString:
		return "OctetString"
	case g.Null:
		return "Null"
	case g.ObjectIdentifier:
		return "ObjectIdentifier"
	case g.IPAddress:
		return "IPAddress"
	case g.Counter32:
		return "Counter32"
	case g.Gauge32:
		return "Gauge32"
	case g.TimeTicks:
		return "TimeTicks"
	case g.Opaque:
		return "Opaque"
	case g.NsapAddress:
		return "NsapAddress"
	case g.Counter64:
		return "Counter64"
	case g.Uinteger32:
		return "Uinteger32"
	case g.NoSuchObject:
		return "NoSuchObject"
	case g.NoSuchInstance:
		return "NoSuchInstance"
	case g.EndOfMibView:
		return "EndOfMibView"
	default:
		return fmt.Sprintf("Unknown(%d)", t)
	}
}
