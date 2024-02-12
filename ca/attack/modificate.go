package attack

import (
	"fmt"
	"log/slog"

	"github.com/waterm310n/rpkiemu-go/ca/krillop"
)

func modificateAttack(attack *Attack, caOps map[string]*krillop.KrillK8sCA) {
	switch attack.AttackObject {
	case AttacKObject_RESOURCECERT:
		//Todo AttacKObject_RESOURCECERT
		fmt.Print("123")
	case AttacKObject_ROAS:
		for _, attackData := range attack.AttackData {
			publishPoint := attackData.PublishPoint
			handle := attackData.HandleName
			for _, binding := range attackData.PreBindings {
				caOps[publishPoint].RemoveAsnIpPair(handle, binding.Ip,fmt.Sprintf("%v",binding.Asn))
			}
			for _, binding := range attackData.AfterBindings {
				caOps[publishPoint].AddAsnIpPair(handle, binding.Ip, fmt.Sprintf("%v",binding.Asn))
			}
		}
	default:
		slog.Error(fmt.Sprintf("AttackObject %v is wrong", attack.AttackObject))
		return
	}
}
