package attack

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/waterm310n/rpkiemu-go/ca/krillop"
)

type AttackType int32

const (
	AttackType_REVOCATE   AttackType = 0
	AttackType_MODIFICATE AttackType = 1
	AttackType_CORRUPT    AttackType = 2
	AttackType_DELETE     AttackType = 3
	AttackType_INJECT     AttackType = 4
)

var (
	AttackType_name = map[int32]string{
		0: "REVOCATE",
		1: "MODIFICATE",
		2: "CORRUPT",
		3: "DELETE",
		4: "INJECT",
	}
	AttackType_value = map[string]int32{
		"REVOCATE":   0,
		"MODIFICATE": 1,
		"CORRUPT":    2,
		"DELETE":     3,
		"INJECT":     4,
	}
)

func (t *AttackType) UnmarshalJSON(bytes []byte) error {
	s := string(bytes[1 : len(bytes)-1])
	if num, ok := AttackType_value[s]; ok {
		*t = AttackType(num)
		return nil
	}
	return fmt.Errorf("cant parse %s while unmarshal AttackType during (t *AttackType) UnmarshalJSON(bytes []byte) error call", bytes)
}

type AttacKObject int32

const (
	AttacKObject_RESOURCECERT AttacKObject = 0
	AttacKObject_ROAS         AttacKObject = 1
)

var (
	AttacKObject_name = map[int32]string{
		0: "RESOURCECERT",
		1: "ROA",
	}
	AttacKObject_value = map[string]int32{
		"RESOURCECERT": 0,
		"ROA":          1,
	}
)

func (t *AttacKObject) UnmarshalJSON(bytes []byte) error {
	s := string(bytes[1 : len(bytes)-1])
	if num, ok := AttacKObject_value[s]; ok {
		*t = AttacKObject(num)
		return nil
	}
	return fmt.Errorf("cant parse %s while unmarshal AttacKObject during (t *AttacKObject) UnmarshalJSON(bytes []byte) error call", bytes)
}

type Binding struct {
	Ip  string `json:"ip,omitempty"`
	Asn int    `json:"asn,omitempty"`
}

type AttackData struct {
	HandleName         string     `json:"handle_name,omitempty"`
	PublishPoint       string     `json:"publish_point,omitempty"`
	ParentPublishPoint string     `json:"parent_publish_point,omitempty"`
	ParentHandleName   string     `json:"parent_handle_name,omitempty"`
	ASN                int        `json:"asn,omitempty"`
	Ipv4Resource       []string   `json:"ipv_4_resource,omitempty"`
	Ipv6Resource       []string   `json:"ipv_6_resource,omitempty"`
	PreBindings        []*Binding `json:"pre_bindings,omitempty"`
	Bindings           []*Binding `json:"bindings,omitempty"`
	AfterBindings      []*Binding `json:"after_bindings,omitempty"`
}

type Attack struct {
	AttackType   AttackType    `json:"attack_type,omitempty"`
	AttackObject AttacKObject  `json:"attack_object,omitempty"`
	AttackData   []*AttackData `json:"attack_data,omitempty"`
}

func parseAttackJson(attackJson string) ([]*Attack, error) {
	var attacks []*Attack
	if content, err := os.ReadFile(attackJson); err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal(content, &attacks); err != nil {
			return nil, err
		}
		return attacks, nil
	}
}

func ExcuteAttack(attackJson string) {
	attacks, err := parseAttackJson(attackJson)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	caOps := krillop.CreateCAOp()
	for _, attack := range attacks {
		switch attack.AttackType {
		case AttackType_REVOCATE:
			revocateAttack(attack, caOps)
		case AttackType_MODIFICATE:
			modificateAttack(attack, caOps)
		case AttackType_CORRUPT:
			corruptAttack(attack, caOps)
		case AttackType_DELETE:
			deleteAttack(attack, caOps)
		case AttackType_INJECT:
			injectAttack(attack, caOps)
		default:
			slog.Error(fmt.Sprintf("AttackType %v is wrong", attack.AttackType))
			break
		}
	}
}
