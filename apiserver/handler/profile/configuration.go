package profile

import (
	"fmt"
	M "probe/model"
	"time"
)

var configurationKeys = []string{
	"token",
	"id",
	"sid",
	"name",
	"ip_address",
	"hostname",
	"osname",
	"role",
	"secret_key",
	"org_id",
	"created",
	"modified",
}

type TokenPayload struct {
	Token string `json:"token"`
}

type TextPayload struct {
	T string `json:"text"`
}

/*
Basic agent info keys - reference from product

// "organization_key": "0euEc91hD2+sX3xkwPMv481ynRZSbGXeAsqjTGA2wkwg=",
// "agent_id": "AB97709CF07245119537BF980B925197",
// "agent_sid": "87D7E472-D6FE-4B02-9C3B-52103699AE7F",
// "agent_ip": "10.240.44.148",
// "agent_name": "AB97709CF07245119537BF980B925197",
// "agent_hostname": "mnvl-ta-everest-dc01p1.india.airtel.itm",
// "agent_osname": "Red Hat Enterprise Linux 8.9 (Ootpa)",
// "primary_agent_id": "",
// "primary_agent_failure_count": 2,
// "agent_role": "primary",
// "registered_in_cloud": 1,
// "permission_level": 8139,
// "agent_type": 1,
// "creation_time": "2024-03-22T10:01:47.598Z",
// "updation_time": "2024-03-22T10:01:47.598Z",
// "agent_passkey": "2123769780694589440",
// "kafka_access_key": "infraon-producer",
// "kafka_secret_key": "infraon@123",
// "org_id": "126399850273257295872"
*/

func (api *R) getAgentConfiguration() (data map[string]any, err error) {
	db, ok := api.DB.Get(string(M.CONFIG))
	if !ok {
		return nil, fmt.Errorf("NO_CONFIGURATION")
	}

	data = make(map[string]any)
	for _, k := range configurationKeys {
		r, _ := db.Find(k)
		data[k] = string(r)
	}

	return
}

func (api *R) updateAgentToken(token string) (err error) {
	var key = "token"
	db, ok := api.DB.Get(string(M.CONFIG))
	if !ok {
		return fmt.Errorf("NO_CONFIGURATION")
	}

	// checking token exists for not
	r, err := db.Find(key)
	if err == nil {
		if len(r) > 1 {
			return fmt.Errorf("AGENT_TOKEN_EXISTS")
		}
	}

	// crating new token into system
	err = db.Create(key, []byte(token))
	if err != nil {
		return err
	}
	_now := time.Now().UTC().String()
	db.Create("created", []byte(_now))
	db.Create("modified", []byte(_now))

	// TODO
	// need to create more keys into db

	return
}
